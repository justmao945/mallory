package mallory

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

// all write on this should be sync between threads
type EngineGAE struct {
	// Global config
	Env *Env
	// Loaded certificate, contains the root certificate and private key
	RootCA *tls.Certificate
	// Pool of auto generated fake certificates signed by RootCert
	Certs *CertPool
}

func NewEngineGAE(e *Env) *EngineGAE {
	return &EngineGAE{Env: e}
}

func (self *EngineGAE) Init() (err error) {
	rcert, err := tls.LoadX509KeyPair(self.Env.Cert, self.Env.Key)
	if err != nil {
		return
	}
	self.RootCA = &rcert
	self.Certs = NewCertPool()
	return
}

// 1. Receive client request R1
// 2. Write R1 as the body of a new request R2
// 3. Post request R2 to remote GAE
// 4. Receive response P1 from GAE
// 5. Read remote server(which the client want to connect with) resonse P2 from the body of P1
// 6. Send P2 as the response to client
func (self *EngineGAE) Serve(s *Session) {
	w, r := s.ResponseWriter, s.Request
	if r.Method == "CONNECT" {
		s.Error("this function can not handle CONNECT method")
		return
	}

	// write the client request and post to remote
	// Note: WriteProxy keeps the full request URI
	var buf bytes.Buffer
	if err := r.WriteProxy(&buf); err != nil {
		s.Error("http.Request.WriteProxy: %s", err.Error())
		return
	}

	// use httpS to keep all things secure,
	// the second phase of CONNECT also uses this.
	url := fmt.Sprintf("https://%s.appspot.com/http", self.Env.AppSpot)
	// for debug
	if self.Env.AppSpot == "debug" {
		url = "http://localhost:8080/http"
	}

	// post client request as body data
	resp, err := http.Post(url, "application/data", &buf)
	if err != nil {
		s.Error("http.Post: %s", err.Error())
		return
	}

	// the response for the requst of client
	cres, err := http.ReadResponse(bufio.NewReader(resp.Body), r)
	if err != nil {
		s.Error("http.ReadResponse: %s", err.Error())
		return
	}

	// copy headers
	CopyResponseHeader(w, cres)

	// please prepare header first and write them
	w.WriteHeader(cres.StatusCode)

	_, err = io.Copy(w, cres.Body)
	if err != nil {
		s.Error("io.Copy: %s", err.Error())
		return
	}

	// Must close body after read the response body
	// Note that cres.Body is rely on resp.Body, so do not close before reading
	if err := resp.Body.Close(); err != nil {
		s.Error("http.Response.Body.Close: %s", err.Error())
		return
	}

	s.Info("RESPONSE %s %s", r.URL.Host, resp.Status)
}

// FIXME: Impossible to connect gae and handle it as a normal TCP connection?
// GAE only provide http handlers? At least I don't know how to handle to TCP connection on GAE server.
// NOTE: GAE socket service can only be available for billing users. So free users is unable to use the
// long term connection. And do what we did in EngineDirect.
// So we can only use urlfetch.Client.Transport.RoundTrip to do http or https method.
// Generally, the CONNECT method can be used for any purpose for the advantage of TCP connection.
// The proxy doesn't need to know what the real underlying protocol or what it is, just need to copy
// data from client to server, and copy the response from the server to client without any interpret.
// Now what we can do and had been done by some GAE proxies is that, extract the underlying protocol!!!
// GAE can only handle limited protocols with urlfetch module, such as http and https.
// Use Hijacker to get the underlying connection
//
// 1. Detect host and port, currectly only support 443 HTTPS request
// 2. Hijack the client connection
// 3. Dial self
// 4. Return 200 OK if is successfully
// 5. Get cached or create new signed certificate
// 6. Wrap client connection with TLS and make handshake
// 7. Receive http request
// 8. Write request as a proxy request to self, HTTP handler
// 9. Copy response to client...
func (self *EngineGAE) Connect(s *Session) {
	w, r := s.ResponseWriter, s.Request
	if r.Method != "CONNECT" {
		s.Error("this function can only handle CONNECT method")
		return
	}

	// Only support HTTPS protocol, which is connected with port 443
	host, port, err := net.SplitHostPort(r.URL.Host)
	if err != nil {
		s.Error("net.SplitHostPort: %s", err.Error())
		return
	}

	if port != "443" {
		s.Error("unsupported CONNECT port: %s", port)
		return
	}

	// hijack the connection to make SSL handshake
	hij, ok := w.(http.Hijacker)
	if !ok {
		s.Error("Server does not support Hijacker")
		return
	}

	conn, _, err := hij.Hijack()
	defer conn.Close()
	if err != nil {
		s.Error("http.Hijacker.Hijack: %s", err.Error())
		return
	}

	// dial self to transport application data, http request
	rconn, err := net.Dial("tcp", self.Env.Addr)
	defer rconn.Close()
	if err != nil {
		s.Error("net.Dial: %s", err.Error())
		return
	}

	// Once connected successfully, return OK
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	// get the fake cert, every host should have its own cert
	cert := self.Certs.GetSafe(host)
	if cert == nil {
		config := &CertConfig{
			SerialNumber: s.ID,
			CommonName:   host, // FIXME: common name mismatch
		}
		cert, err = CreateSignedCert(self.RootCA, config)
		if err != nil {
			s.Error("EngineGAE.CreateSignedCert: %s", err.Error())
			return
		}
		self.Certs.AddSafe(host, cert)
	}

	// assume the protocol of client connection is HTTPS
	// wrap it with TSL server
	config := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ServerName:   host,
	}
	sconn := tls.Server(conn, config)
	defer sconn.Close()

	// The TLS connection goes here
	if err := sconn.Handshake(); err != nil {
		// re-open browser to recover the handshake error:
		//    remote error: bad certificate
		s.Error("tls.Server.Handshake: %s", err.Error())
		return
	}

	// finally, we are at application layer, http request comes
	req, err := http.ReadRequest(bufio.NewReader(sconn))
	if err != nil {
		s.Error("http.ReadRequest: %s", err.Error())
		return
	}

	// should re-wrap the URL with scheme "https://"
	req.URL, err = url.Parse("https://" + r.Host + req.URL.String())

	// Now re-write the client request to self, HTTP handler
	err = req.WriteProxy(rconn)
	if err != nil {
		s.Error("http.Request.WriteProxy: %s", err.Error())
		return
	}

	// copy response
	_, err = io.Copy(sconn, rconn)
	if err != nil {
		s.Error("io.Copy: %s", err.Error())
		return
	}

	s.Info("CLOSE %s", r.URL.Host)
}
