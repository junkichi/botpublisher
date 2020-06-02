package main

import (
	"fmt"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"go.mongodb.org/mongo-driver/mongo"
)

var urayasuTagTweetDbID = "test"
var urayasuTagTweetDbCOL = "tweet"

func initUrayasuTagTweet(c mongo.Client, api *anaconda.TwitterApi, query string) {
	col := c.Database(urayasuTagTweetDbID).Collection(urayasuTagTweetDbCOL)
	var tweets []anaconda.Tweet
	TWSearchTweet(api, query, &tweets)

	n := 0
	skip := 0
	for _, tweet := range tweets {
		found, err := DBFindTweet(col, tweet.User.ScreenName, tweet.IdStr)
		if found == true {
			skip++
			continue
		}
		link := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IdStr)
		pubdate, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
		err = DBInsertTweet(col,
			tweet.User.Name,
			tweet.User.ScreenName,
			tweet.IdStr,
			link,
			tweet.FullText,
			pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		n++
	}
	fmt.Println("[utag] skipped:", skip)
	fmt.Println("[utag] inserted:", n)
}

func collectUrayasuTagTweet(c mongo.Client, api *anaconda.TwitterApi, query string) {
	pcol := c.Database(publishDbID).Collection(publishDbCOL)
	col := c.Database(urayasuTagTweetDbID).Collection(urayasuTagTweetDbCOL)
	var tweets []anaconda.Tweet
	TWSearchTweet(api, query, &tweets)

	n := 0
	skip := 0
	for _, tweet := range tweets {
		found, err := DBFindTweet(col, tweet.User.ScreenName, tweet.IdStr)
		if found == true {
			skip++
			continue
		}
		link := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IdStr)
		pubdate, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
		err = DBInsertTweet(col,
			tweet.User.Name,
			tweet.User.ScreenName,
			tweet.IdStr,
			link,
			tweet.FullText,
			pubdate)
		if err != nil {
			fmt.Println("insert error:", err)
			continue
		}
		desc := fmt.Sprintln("(浦安タグ)", tweet.User.Name, link)
		err = DBInsertPublish(pcol, desc)
		if err != nil {
			fmt.Println("insert error:", err)
		}
		n++
	}
	fmt.Println("[utag] skipped:", skip)
	fmt.Println("[utag] inserted:", n)
}
