package go_sitemap_parser

import (
	"testing"
	"os"
	"fmt"
)

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

func TestSitemapParse3(t *testing.T) {
	urlsChan := make(chan SitemapUrl, 1)

	go func() {
		err := (&Sitemap{"https://www.ozon.ru/sitemap.xml"}).Parse(urlsChan)

		if err != nil {
			t.Fatal(err)
		}
	}()

	file, _ := os.Create("log")
	for url := range urlsChan {
		file.WriteString(fmt.Sprintln(url))
	}
}