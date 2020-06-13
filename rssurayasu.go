package main

import (
	"botpublisher/rss"
	"botpublisher/storage"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

var urayasuRssURL = "http://www.city.urayasu.lg.jp/news.rss"
var urayasuRssCOL = "rss"

type UrayasuRSS struct{}

func (UrayasuRSS) Init() bool {
	s := storage.GetInstance()
	var items []*gofeed.Item
	err := rss.RetrieveRSS(urayasuRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	for _, item := range items {
		found, err := storage.FindRSS(s, urayasuRssCOL, item.Link)
		if found == true {
			skip++
			continue
		}
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.Published)
		err = storage.InsertRSS(s, urayasuRssCOL, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		n++
	}
	fmt.Println("[urss] skipped:", skip)
	fmt.Println("[urss] inserted:", n)

	return true
}

func (UrayasuRSS) Collect() bool {
	s := storage.GetInstance()
	var items []*gofeed.Item
	err := rss.RetrieveRSS(urayasuRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	for _, item := range items {
		found, err := storage.FindRSS(s, urayasuRssCOL, item.Link)
		if found == true {
			skip++
			continue
		}
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.Published)
		err = storage.InsertRSS(s, urayasuRssCOL, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		desc := fmt.Sprintln("(浦安市)", item.Title, item.Link)
		err = storage.InsertPublish(s, publishCOL, desc)
		if err != nil {
			fmt.Println("insert error:", err)
		}
		n++
	}
	fmt.Println("[urss] skipped:", skip)
	fmt.Println("[urss] inserted:", n)

	return true
}
