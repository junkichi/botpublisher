package main

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/mongo"
)

var urayasuRssURL = "http://www.city.urayasu.lg.jp/news.rss"
var urayasuRssDbID = "test"
var urayasuRssDbCOL = "rss"

func initUrayasuRSS(c mongo.Client) bool {
	col := c.Database(urayasuRssDbID).Collection(urayasuRssDbCOL)
	var items []*gofeed.Item
	err := RetrieveRSS(urayasuRssURL, &items)
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
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.Published)
		err = DBInsertRSS(col, item.Title, item.Link, item.Description, pubdate)
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

func collectUrayasuRSS(c mongo.Client) bool {
	pcol := c.Database(publishDbID).Collection(publishDbCOL)
	col := c.Database(urayasuRssDbID).Collection(urayasuRssDbCOL)
	var items []*gofeed.Item
	err := RetrieveRSS(urayasuRssURL, &items)
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
		pubdate, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.Published)
		err = DBInsertRSS(col, item.Title, item.Link, item.Description, pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		desc := fmt.Sprintln("(浦安市)", item.Title, item.Link)
		err = DBInsertPublish(pcol, desc)
		if err != nil {
			fmt.Println("insert error:", err)
		}
		n++
	}
	fmt.Println("[urss] skipped:", skip)
	fmt.Println("[urss] inserted:", n)

	return true
}
