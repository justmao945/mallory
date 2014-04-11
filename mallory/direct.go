package mallory

import (
	"io"
	"net"
	"net/http"
	"sync"
)

type EngineDirect struct{}

func NewEngineDirect(e *Env) *EngineDirect {
	return &EngineDirect{}
}

func (self *EngineDirect) Init() (err error) {
	return
}

// 1. Receive request R1 from client
// 2. Re-post request R1 to remote server(the one client want to connect)
// 3. Receive response P1 from remote server
// 4. Send response P1 to client
func (self *EngineDirect) Serve(s *Session) {
	w, r := s.ResponseWriter, s.Request
	if r.Method == "CONNECT" {
		s.Error("this function can not handle CONNECT method")
		return
	}

	// Client.Do is different from DefaultTransport.RoundTrip ...
	// Client.Do will change some part of request as a new request of the server.
	// The underlying RoundTrip never changes anything of the request.
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		s.Error("http.DefaultTransport.RoundTrip: %s", err.Error())
		return
	}

	// copy headers
	CopyResponseHeader(w, resp)

	// please prepare header first and write them
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		s.Error("io.Copy: %s", err.Error())
		return
	}

	// Must close body after read the response body
	if err := resp.Body.Close(); err != nil {
		s.Error("http.Response.Body.Close: %s", err.Error())
		return
	}
	s.Info("RESPONSE %s %s", r.URL.Host, resp.Status)
}

// 1. Receive CONNECT request from the client
// 2. Dial the remote server(the one client want to conenct)
// 3. Send 200 OK to client if the connection is established
// 4. Exchange data between client and server
func (self *EngineDirect) Connect(s *Session) {
	w, r := s.ResponseWriter, s.Request
	if r.Method != "CONNECT" {
		s.Error("this function can only handle CONNECT method")
		return
	}

	// Use Hijacker to get the underlying connection
	hij, ok := w.(http.Hijacker)
	if !ok {
		s.Error("Server does not support Hijacker")
		return
	}

	src, _, err := hij.Hijack()
	defer src.Close()
	if err != nil {
		s.Error("http.Hijacker.Hijack: %s", err.Error())
		return
	}

	// connect the remote client directly
	dst, err := net.Dial("tcp", r.URL.Host)
	defer dst.Close()
	if err != nil {
		s.Error("net.Dial: %s", err.Error())
		return
	}

	// Once connected successfully, return OK
	src.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	// Proxy is no need to know anything, just exchange data between the client
	// the the remote server.
	var wg sync.WaitGroup
	wg.Add(2)

	copyAndWait := func(w io.Writer, r io.Reader) {
		_, err := io.Copy(w, r)
		if err != nil {
			s.Error("io.Copy: %s", err.Error())
		}
		wg.Done()
	}
	go copyAndWait(dst, src)
	go copyAndWait(src, dst)

	// Generally, the remote server would keep the connection alive,
	// so we will not close the connection until both connection recv
	// EOF and are done!
	wg.Wait()

	s.Info("CLOSE %s", r.URL.Host)
}
