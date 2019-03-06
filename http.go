package mallory

import (
	"fmt"
	"net"
	"net/http"
)

// HostOnly returns host if has port in addr, or addr if missing port
func HostOnly(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	} else {
		return host
	}
}

// copy and overwrite headers from r to w
func CopyHeader(w http.ResponseWriter, r *http.Response) {
	// copy headers
	dst, src := w.Header(), r.Header
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

// StatusText returns http status text looks like "200 OK"
func StatusText(c int) string {
	return fmt.Sprintf("%d %s", c, http.StatusText(c))
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	// If no Accept-Encoding header exists, Transport will add the headers it can accept
	// and would wrap the response body with the relevant reader.
	"Accept-Encoding",
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
	"Proxy-Connection", // added by CURL  http://homepage.ntlworld.com/jonathan.deboynepollard/FGA/web-proxy-connection-header.html
}

func RemoveHopHeaders(h http.Header) {
	for _, k := range hopHeaders {
		h.Del(k)
	}
}
