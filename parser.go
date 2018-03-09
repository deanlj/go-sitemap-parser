package go_sitemap_parser

import (
	"errors"
	"net/http"
	"bufio"
	"io"
	"regexp"
	"strings"
	"sync"
)

var tagRegexp = regexp.MustCompile(`^\s*<(/?\w+)`)

type SitemapUrl struct {
	Loc 		string
	Lastmod 	string
	Changefreq 	string
	Priority 	string
}

type Sitemap struct {
	Loc string
}

func (s *Sitemap) Parse(urlsChan chan SitemapUrl) error {
	waitGroup := &sync.WaitGroup{}
	sitemaps := &sync.Map{}
	err := s.parse(urlsChan, waitGroup, sitemaps)

	waitGroup.Wait()
	close(urlsChan)

	return err
}

func (s *Sitemap) parse(urlsChan chan SitemapUrl, waitGroup *sync.WaitGroup, sitemaps *sync.Map) error {
	if s.Loc == "" {
		return errors.New("sitemap loc is not defined")
	}

	sitemapResp, err := http.Get(s.Loc)

	if err != nil {
		return err;
	}

	return (&sitemapParser{
		UrlsChan: 	urlsChan,
		Sitemaps: 	sitemaps,
		WaitGroup: 	waitGroup,
	}).parseResp(sitemapResp)
}

type sitemapParser struct {
	UrlsChan 	chan SitemapUrl
	Sitemaps	*sync.Map
	WaitGroup 	*sync.WaitGroup

	sitemap 	bool
	url 		*SitemapUrl
	loc 		bool
	lastmod 	bool
	changefreq 	bool
	priority 	bool

	err        error
}

func (p *sitemapParser) parseResp(resp *http.Response) error {
	bufRespReader := bufio.NewReader(resp.Body)
	p.err = nil

	var str string
	var err error

	for ; ; str, err = bufRespReader.ReadString('>') {
		p.parseStr(str)

		if err != nil || p.err != nil {
			break
		}
	}

	if p.err != nil {
		err = p.err
	}

	if err == io.EOF {
		err = nil
	}

	return err
}

func (p *sitemapParser) parseStr(str string) {
	if !p.parseTag(str) {
		switch true {
			case p.sitemap:    p.parseSitemap(str)
			case p.url != nil: p.parseUrl(str)
		}
	}
}

func (p *sitemapParser) parseTag(str string) bool {
	submatch := tagRegexp.FindStringSubmatch(str)

	if len(submatch) < 2 {
		return false
	}

	tag := submatch[1]

	if tag == "" {
		return false
	}

	switch tag {
		case "sitemap": 	p.sitemap 		= true
		case "/sitemap": 	p.sitemap 		= false
		case "url": 		p.url 			= &SitemapUrl{}
		case "/url": 		p.pushUrl()
		case "loc": 		p.loc 			= true
		case "lastmod": 	p.lastmod 		= true
		case "changefreq": 	p.changefreq 	= true
		case "priority": 	p.priority 		= true
	}

	return true
}

func (p *sitemapParser) pushUrl() {
	p.UrlsChan <- *p.url
	p.url = nil
}

func (p *sitemapParser) parseSitemap(str string) {
	if p.loc {
		index := strings.Index(str, "</loc>")

		if index > 0 {
			loc := str[0:index]

			if _, loaded := p.Sitemaps.LoadOrStore(loc, true); !loaded {
				p.WaitGroup.Add(1)

				go func() {
					p.err = (&Sitemap{loc}).parse(p.UrlsChan, p.WaitGroup, p.Sitemaps)
					p.WaitGroup.Done()
				}()
			}
		}
	}

	p.loc 			= false
	p.lastmod 		= false
	p.changefreq 	= false
	p.priority 		= false
}

func (p *sitemapParser) parseUrl(str string) {
	switch true {
		case p.loc:
			index := strings.Index(str, "</loc>")
			if index > 0 {
				p.url.Loc = str[0:index]
			}

			p.loc = false

		case p.lastmod:
			index := strings.Index(str, "</lastmod>")
			if index > 0 {
				p.url.Lastmod = str[0:index]
			}

			p.lastmod = false

		case p.changefreq:
			index := strings.Index(str, "</changefreq>")
			if index > 0 {
				p.url.Changefreq = str[0:index]
			}

			p.changefreq = false

		case p.priority:
			index := strings.Index(str, "</priority>")
			if index > 0 {
				p.url.Priority = str[0:index]
			}

			p.priority = false
	}
}