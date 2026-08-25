package httpapi

import (
	"net/http"
	"strings"
)

func bearer(r *http.Request) (string, bool) {
	v := r.Header.Get("Authorization")
	if v == "" {
		return "", false
	}
	p := strings.SplitN(v, " ", 2)
	if len(p) != 2 || !strings.EqualFold(p[0], "Bearer") || p[1] == "" {
		return "", false
	}
	return p[1], true
}
func setNoCache(w http.ResponseWriter) { w.Header().Set("Cache-Control", "no-store") }
