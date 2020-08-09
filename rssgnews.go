package main

import (
	"fmt"
	"time"

	"github.com/junkichi/botpublisher/rss"
	"github.com/junkichi/botpublisher/storage"

	"github.com/mmcdole/gofeed"
)

var googleNewsRssURL = "https://news.google.com/rss/search?q=%E6%B5%A6%E5%AE%89&hl=ja&gl=JP&ceid=JP:ja"
var googleNewsCOL = "rss"

// GoogleNewsRSS is the value of iamge directory
type GoogleNewsRSS struct {
	imgdir string
}

// GetImageDir the stored value
func (r GoogleNewsRSS) GetImageDir() string {
	return r.imgdir
}

// Init the RSS
func (GoogleNewsRSS) Init() bool {
	s := storage.GetInstance()
	var items []*gofeed.Item
	err := rss.RetrieveRSS(googleNewsRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	for _, item := range items {
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.Published)
		found, err := storage.FindRSS(s, googleNewsCOL, item.Link, pubdate)
		if found == true {
			skip++
			continue
		}
		err = storage.InsertRSS(s, googleNewsCOL, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		n++
	}
	fmt.Println("[gnews] skipped:", skip)
	fmt.Println("[gnews] inserted:", n)

	return true
}

// Collect the RSS
func (GoogleNewsRSS) Collect(imgdir string) bool {
	s := storage.GetInstance()
	var items []*gofeed.Item
	err := rss.RetrieveRSS(googleNewsRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	old := 0
	for _, item := range items {
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.Published)
		found, err := storage.FindRSS(s, googleNewsCOL, item.Link, pubdate)
		if found == true {
			skip++
			continue
		}
		err = storage.InsertRSS(s, googleNewsCOL, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		limitdate := time.Now().Add(-12 * time.Hour)
		if pubdate.Before(limitdate) {
			old++
			continue
		}
		desc := fmt.Sprintln("(News)", item.Title, item.Link)
		err = storage.InsertPublish(s, publishCOL, desc, "")
		if err != nil {
			fmt.Println("insert error:", err)
		}
		n++
	}
	fmt.Println("[gnews] skipped:", skip)
	fmt.Println("[gnews] old:", old)
	fmt.Println("[gnews] inserted:", n)

	return true
}
