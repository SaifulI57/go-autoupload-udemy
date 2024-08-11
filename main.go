package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/alitto/pond"
)

func main() {
	// messager.Start()
	f, err := os.Open("test.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	pool := pond.New(0, 100)
	var photos []string
	databefore, _ := doc.Find("a.tme_messages_more").Attr("data-before")
	doc.Find("a.tgme_widget_message_photo_wrap").Each(func(i int, s *goquery.Selection) {
		pool.Submit(func() {
			href, _ := s.Attr("style")
			re := regexp.MustCompile(`background-image:url\(['\"]?(.*?)['\"]?\)`)
			match := re.FindStringSubmatch(href)
			if len(match) > 1 {
				photos = append(photos, match[1])
			}
		})
	})
	pool.StopAndWait()
	type CourseInfo struct {
		Title            string
		DescriptionText  string
		DescriptionValue string
		CategoryText     string
		CategoryValue    string
		LengthText       string
		LengthValue      string
		EnrollLinkText   string
		EnrollLinkValue  string
		PhotoLink        string
	}

	var messages []CourseInfo
	temps := CourseInfo{}
	lphoto := 0
	doc.Find("div.tgme_widget_message_text").Contents().Each(func(i int, s *goquery.Selection) {

		switch {
		case i%36 == 0:
			temps.Title = s.Text()
		case i%36 == 4:
			temps.DescriptionText = s.Text()
		case i%36 == 5:
			temps.DescriptionValue = s.Text()
		case i%36 == 9:
			temps.CategoryText = s.Text()
		case i%36 == 10:
			temps.CategoryValue = s.Text()
		case i%36 == 14:
			temps.LengthText = s.Text()
		case i%36 == 15:
			temps.LengthValue = s.Text()
		case i%36 == 19:
			temps.EnrollLinkText = s.Text()
		case i%36 == 21:
			temps.EnrollLinkValue = s.Text()
		case i%36 == 35:
			temps.PhotoLink = photos[lphoto]
			lphoto++
			messages = append(messages, temps)
			temps = CourseInfo{}
		}

	})

	fmt.Println(databefore)
	// fmt.Println(photos)
	fmt.Println(messages)

}
