package main

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/mongo"
)

var googleNewsRssURL = "https://news.google.com/rss/search?q=%E6%B5%A6%E5%AE%89&hl=ja&gl=JP&ceid=JP:ja"
var googleNewsDbID = "test"
var googleNewsDbCOL = "rss"

func initGoogleNewsRSS(c mongo.Client) bool {
	col := c.Database(googleNewsDbID).Collection(googleNewsDbCOL)
	var items []*gofeed.Item
	err := RetrieveRSS(googleNewsRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	for _, item := range items {
		found, err := DBFindRSS(col, item.Link)
		if found == true {
			skip++
			continue
		}
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.Published)
		err = DBInsertRSS(col, item.Title, item.Link, item.Description, pubdate)
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

func collectGoogleNewsRSS(c mongo.Client) bool {
	pcol := c.Database(publishDbID).Collection(publishDbCOL)
	col := c.Database(googleNewsDbID).Collection(googleNewsDbCOL)
	var items []*gofeed.Item
	err := RetrieveRSS(googleNewsRssURL, &items)
	if err != nil {
		fmt.Println("retrieve err:", err)
		return false
	}

	n := 0
	skip := 0
	for _, item := range items {
		found, err := DBFindRSS(col, item.Link)
		if found == true {
			skip++
			continue
		}
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 MST", item.Published)
		err = DBInsertRSS(col, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		desc := fmt.Sprintln("(News)", item.Title, item.Link)
		err = DBInsertPublish(pcol, desc)
		if err != nil {
			fmt.Println("insert error:", err)
		}
		n++
	}
	fmt.Println("[gnews] skipped:", skip)
	fmt.Println("[gnews] inserted:", n)

	return true
}
