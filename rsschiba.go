package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/junkichi/botpublisher/rss"
	"github.com/junkichi/botpublisher/storage"

	"github.com/mmcdole/gofeed"
)

var chibaRssURL = "https://www.pref.chiba.lg.jp/homepage/shinchaku/shinchaku.xml"
var chibaRssCOL = "rss"

// ChibaRSS is the value of image directory
type ChibaRSS struct {
	imgdir string
}

// GetImageDir the stored balue
func (r ChibaRSS) GetImageDir() string {
	return r.imgdir
}

// Init the RSS
func (ChibaRSS) Init() bool {
	s := storage.GetInstance()
	var items []*gofeed.Item
	err := rss.RetrieveRSS(chibaRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	for _, item := range items {
		pubdate, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", item.Published)
		found, err := storage.FindRSS(s, chibaRssCOL, item.Link, pubdate)
		if found == true {
			skip++
			continue
		}
		found = strings.Contains(item.Title, "新型コロナウイルス")
		if found == false {
			skip++
			continue
		}
		err = storage.InsertRSS(s, chibaRssCOL, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		n++
	}
	fmt.Println("[crss] skipped:", skip)
	fmt.Println("[crss] inserted:", n)

	return true
}

// Collect the RSS
func (ChibaRSS) Collect(imgdir string) bool {
	s := storage.GetInstance()
	var items []*gofeed.Item
	err := rss.RetrieveRSS(chibaRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	old := 0
	for _, item := range items {
		pubdate, _ := time.Parse("2006-01-02T15:04:05Z07:00", item.Published)
		found, err := storage.FindRSS(s, chibaRssCOL, item.Link, pubdate)
		if found == true {
			skip++
			continue
		}
		found = strings.Contains(item.Title, "新型コロナウイルス")
		if found == false {
			skip++
			continue
		}
		err = storage.InsertRSS(s, chibaRssCOL, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		imgpath := fmt.Sprintf("%s/%d.png", imgdir, time.Now().Unix())
		err = rss.TakeScreenShot(item.Link, `#tmp_contents`, imgpath)
		if err != nil {
			imgpath = ""
		}
		limitdate := time.Now().Add(-24 * 7 * time.Hour)
		if pubdate.Before(limitdate) {
			old++
			continue
		}
		desc := fmt.Sprintln("(千葉県)", item.Title, item.Link)
		err = storage.InsertPublish(s, publishCOL, desc, imgpath)
		if err != nil {
			fmt.Println("insert error:", err)
		}
		n++
	}
	fmt.Println("[crss] skipped:", skip)
	fmt.Println("[crss] old:", old)
	fmt.Println("[crss] inserted:", n)

	return true
}
