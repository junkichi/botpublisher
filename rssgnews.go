package main

import (
	"botpublisher/storage"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

var googleNewsRssURL = "https://news.google.com/rss/search?q=%E6%B5%A6%E5%AE%89&hl=ja&gl=JP&ceid=JP:ja"
var googleNewsCOL = "rss"

func initGoogleNewsRSS() bool {
	s := storage.GetInstance()
	var items []*gofeed.Item
	err := RetrieveRSS(googleNewsRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	for _, item := range items {
		found, err := storage.FindRSS(s, googleNewsCOL, item.Link)
		if found == true {
			skip++
			continue
		}
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.Published)
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

func collectGoogleNewsRSS() bool {
	s := storage.GetInstance()
	var items []*gofeed.Item
	err := RetrieveRSS(googleNewsRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	for _, item := range items {
		found, err := storage.FindRSS(s, googleNewsCOL, item.Link)
		if found == true {
			skip++
			continue
		}
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.Published)
		err = storage.InsertRSS(s, googleNewsCOL, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		desc := fmt.Sprintln("(News)", item.Title, item.Link)
		err = storage.InsertPublish(s, publishCOL, desc)
		if err != nil {
			fmt.Println("insert error:", err)
		}
		n++
	}
	fmt.Println("[gnews] skipped:", skip)
	fmt.Println("[gnews] inserted:", n)

	return true
}
