package query

import (
	"net/url"
	"testing"
)

func TestParseDefaults(t *testing.T) {
	p := Parse(url.Values{})
	if p.Limit != 20 || p.Offset != 0 {
		t.Fatal(p)
	}
}
func TestParseBounds(t *testing.T) {
	p := Parse(url.Values{"limit": []string{"101"}, "offset": []string{"-1"}, "field": []string{"信息技术"}})
	if p.Limit != 20 || p.Offset != 0 || p.Filter != "信息技术" {
		t.Fatal(p)
	}
	p = Parse(url.Values{"limit": []string{"50"}, "offset": []string{"10"}})
	if p.Limit != 50 || p.Offset != 10 {
		t.Fatal(p)
	}
}
func TestParseBadNumbers(t *testing.T) {
	p := Parse(url.Values{"limit": []string{"x"}, "offset": []string{"y"}})
	if p.Limit != 20 || p.Offset != 0 {
		t.Fatal(p)
	}
}
