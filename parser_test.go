package go_sitemap_parser

import "testing"

func TestSitemapParse1(t *testing.T) {
	err := (&Sitemap{"ololo"}).Parse(make(chan SitemapUrl))

	if err == nil {
		t.Fail()
	}
}

func TestSitemapParse2(t *testing.T) {
	err := (&Sitemap{}).Parse(make(chan SitemapUrl))

	if err == nil {
		t.Fail()
	}
}