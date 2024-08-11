package main

import (
	"fmt"
	"regexp"

	"github.com/alitto/pond"
	"github.com/gocolly/colly/v2"
)

type Doscrape struct {
}
type CourseInfo struct {
	Title       string
	Description string
	Category    string
	Length      string
	EnrollLink  string
	PhotoLink   *string
}

var GColly *colly.Collector

func init() {
	GColly = colly.NewCollector(colly.AllowedDomains("t.me"), colly.AllowURLRevisit(), colly.Async(true))
	GColly.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "t.me")
		r.Headers.Set("Sec-Ch-Ua", `"Not:A-Brand";v="99", "Chromium";v="112"`)
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
		r.Headers.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.5615.50 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Pragma", "no-cache")
		r.Headers.Set("Sec-Ch-Ua-Arch", "x86")
		r.Headers.Set("Sec-Ch-Ua-Full-Version", "112.0.5615.50")
	})
}

func (d *Doscrape) ListToday() ([]*CourseInfo, error) {
	pool := pond.New(1, 1000)
	descs := make(chan *CourseInfo, 1000)
	photos := make(chan string, 1000)
	baseUri := "https://t.me/s/udemycoursesfree"
	uri := baseUri

	for i := 0; i < 3; i++ {
		pool.Submit(func() {
			GColly.Visit(uri)
			GColly.OnHTML("div.tgme_container", func(e *colly.HTMLElement) {
				qbefore := e.DOM.ChildrenFiltered("a.tme_messages_more")
				before, _ := qbefore.Attr("data-before")
				uri = fmt.Sprintf("%s?before=%s", baseUri, before)
				qphoto := e.DOM.ChildrenFiltered("a.tgme_widget_message_photo_wrap")
				photo, _ := qphoto.Attr("style")
				re := regexp.MustCompile(`background-image:url\(['\"]?(.*?)['\"]?\)`)
				match := re.FindStringSubmatch(photo)
				var photolink string
				if len(match) > 1 {
					photolink = string(match[1])
				}
				descs <- &CourseInfo{
					Title:       "",
					Description: "",
					Category:    "",
					Length:      "",
					EnrollLink:  "",
					PhotoLink:   &photolink,
				}
			})
			GColly.OnHTML("div.tgme_widget_message_photo_wrap", func(e *colly.HTMLElement) {
				re := regexp.MustCompile(`background-image:url\(['\"]?(.*?)['\"]?\)`)
				style := e.Attr("style")
				match := re.FindStringSubmatch(style)
				if len(match) > 1 {
					photos <- string(match[1])
				}
			})
			GColly.OnHTML("a.tme_messages_more", func(e *colly.HTMLElement) {
				before := e.Attr("data-before")
				uri = fmt.Sprintf("%s?before=%s", baseUri, before)
			})
			GColly.Wait()
		})
	}

	pool.StopAndWait()
	close(descs)
	var desc []*CourseInfo
	var photo []string

	for res := range photos {
		photo = append(photo, res)
	}
	for result := range descs {
		desc = append(desc, result)
	}

	for i, res := range desc {
		res.PhotoLink = &photo[i]
		desc[i] = res
	}

	return desc, nil
}

func (d *Doscrape) Stop() error {

	return nil
}
