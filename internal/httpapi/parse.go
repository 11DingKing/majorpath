package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func decodeBody(r *http.Request, dst any) error {
	if r.Body == nil {
		return io.EOF
	}
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
func contentTypeJSON(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
}
func methodAllowed(w http.ResponseWriter, method string) bool {
	if method == "" {
		return true
	}
	if w.Header().Get("Allow") == "" {
		w.Header().Set("Allow", method)
	}
	return false
}
