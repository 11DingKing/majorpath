package query

import (
	"net/url"
	"strconv"
)

type Page struct {
	Limit, Offset int
	Filter        string
}

func Parse(v url.Values) Page {
	p := Page{Limit: 20}
	if n, _ := strconv.Atoi(v.Get("limit")); n > 0 && n <= 100 {
		p.Limit = n
	}
	if n, _ := strconv.Atoi(v.Get("offset")); n >= 0 {
		p.Offset = n
	}
	p.Filter = v.Get("field")
	return p
}
