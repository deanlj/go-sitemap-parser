package go_sitemap_parser

import (
	"errors"
	"net/http"
	"bufio"
	"io"
	"fmt"
)

type Sitemap struct {
	Loc string
}

func (s *Sitemap) Parse(urlsChan chan SitemapUrl) error {
	if s.Loc == "" {
		return errors.New("sitemap loc is not defined")
	}

	sitemapResp, err := http.Get(s.Loc)

	if err != nil {
		return err;
	}

	return parseResp(sitemapResp, urlsChan)
}

type SitemapUrl struct {
	Loc string
	Lastmod string
	Changefreq string
	Priority string
}


func parseResp(resp *http.Response, urlsChan chan SitemapUrl) error {
	bufRespReader := bufio.NewReader(resp.Body)

	var str string
	var err error

	for ; ; str, err = bufRespReader.ReadString('>') {
		parseErr := parseStr(str, urlsChan)

		if parseErr != nil {
			err = parseErr
		}

		if err != nil {
			break
		}
	}

	if err == io.EOF {
		err = nil
	}

	return err
}

func parseStr(str string, urlChan chan SitemapUrl) error {
	fmt.Println(str)

	return nil
}