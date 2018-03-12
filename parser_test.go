package go_sitemap_parser

import (
	"testing"
	"os"
	"fmt"
	"time"
	"runtime"
	"math"
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
	start := time.Now()
	var maxMem float64 = 0

	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			maxMem = math.Max(float64(m.Alloc), maxMem)
			time.Sleep(100 * time.Millisecond)
		}
	}()

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

	fmt.Println("script time", time.Since(start))
	fmt.Println("script mem", maxMem)
}