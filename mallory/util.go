package mallory

import (
	"net/http"
)

func CopyResponseHeader(w http.ResponseWriter, r *http.Response) {
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
