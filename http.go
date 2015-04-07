package mallory

import (
	"fmt"
	"net/http"
)

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

// Return http status text looks like "200 OK"
func StatusText(c int) string {
	return fmt.Sprintf("%d %s", c, http.StatusText(c))
}
