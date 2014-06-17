package mallory

import (
	"github.com/justmao945/mallory/proxy"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Direct fetcher from the host of proxy
type EngineSocksToHttp struct {
	Env   *Env
	Proxy proxy.Dialer
}

// Create and initialize
func CreateEngineSocksToHttp(e *Env) (self *EngineSocksToHttp, err error) {
	proxyURL, err := url.Parse(e.SocksProxy)
	if err != nil {
		return
	}
	proxyDialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return
	}
	self = &EngineSocksToHttp{
		Env:   e,
		Proxy: proxyDialer,
	}
	return
}

// Data flow:
//  1. Receive request R1 from client
//  2. Re-post request R1 to remote server(the one client want to connect)
//  3. Receive response P1 from remote server
//  4. Send response P1 to client
func (self *EngineSocksToHttp) Serve(s *Session) {
	w, r := s.ResponseWriter, s.Request
	if r.Method == "CONNECT" {
		s.Error("this function can not handle CONNECT method")
		return
	}
	start := time.Now()

	// Client.Do is different from DefaultTransport.RoundTrip ...
	// Client.Do will change some part of request as a new request of the server.
	// The underlying RoundTrip never changes anything of the request.
	tr := http.Transport{Dial: self.Proxy.Dial}
	resp, err := tr.RoundTrip(r)
	if err != nil {
		s.Error("RoundTrip: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	// please prepare header first and write them
	CopyHeader(w, resp)
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		s.Error("Copy: %s", err.Error())
		return
	}

	d := BeautifyDuration(time.Since(start))
	s.Info("RESPONSE %s %s", resp.Status, d)
}

// Data flow:
//  1. Receive CONNECT request from the client
//  2. Dial the remote server(the one client want to conenct)
//  3. Send 200 OK to client if the connection is established
//  4. Exchange data between client and server
func (self *EngineSocksToHttp) Connect(s *Session) {
	w, r := s.ResponseWriter, s.Request
	if r.Method != "CONNECT" {
		s.Error("this function can only handle CONNECT method")
		return
	}
	start := time.Now()

	// Use Hijacker to get the underlying connection
	hij, ok := w.(http.Hijacker)
	if !ok {
		s.Error("Server does not support Hijacker")
		return
	}

	src, _, err := hij.Hijack()
	if err != nil {
		s.Error("Hijack: %s", err.Error())
		return
	}
	defer src.Close()

	// connect the remote client directly
	dst, err := self.Proxy.Dial("tcp", r.URL.Host)
	if err != nil {
		s.Error("Dial: %s", err.Error())
		return
	}
	defer dst.Close()

	// Once connected successfully, return OK
	src.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	// Proxy is no need to know anything, just exchange data between the client
	// the the remote server.
	var wg sync.WaitGroup
	wg.Add(2)

	copyAndWait := func(w io.Writer, r io.Reader) {
		_, err := io.Copy(w, r)
		if err != nil {
			s.Error("Copy: %s", err.Error())
		}
		wg.Done()
	}
	go copyAndWait(dst, src)
	go copyAndWait(src, dst)

	// Generally, the remote server would keep the connection alive,
	// so we will not close the connection until both connection recv
	// EOF and are done!
	wg.Wait()

	s.Info("CLOSE %s", BeautifyDuration(time.Since(start)))
}
