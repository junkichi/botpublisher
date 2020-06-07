package main

import (
	"botpublisher/storage"
	"fmt"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var urayasuTagTweetCOL = "tweet"

func initUrayasuTagTweet(api *anaconda.TwitterApi, query string) {
	s := storage.GetInstance()

	var tweets []anaconda.Tweet
	TWSearchTweet(api, query, &tweets)

	n := 0
	skip := 0
	for _, tweet := range tweets {
		found, err := storage.FindTweet(s, urayasuTagTweetCOL, tweet.User.ScreenName, tweet.IdStr)
		if found == true {
			skip++
			continue
		}
		link := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IdStr)
		pubdate, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
		err = storage.InsertTweet(s, urayasuTagTweetCOL,
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

func collectUrayasuTagTweet(api *anaconda.TwitterApi, query string) {
	s := storage.GetInstance()
	var tweets []anaconda.Tweet
	TWSearchTweet(api, query, &tweets)

	n := 0
	skip := 0
	for _, tweet := range tweets {
		found, err := storage.FindTweet(s, urayasuTagTweetCOL, tweet.User.ScreenName, tweet.IdStr)
		if found == true {
			skip++
			continue
		}
		link := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IdStr)
		pubdate, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
		err = storage.InsertTweet(s, urayasuTagTweetCOL,
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
		err = storage.InsertPublish(s, publishCOL, desc)
		if err != nil {
			fmt.Println("insert error:", err)
		}
		n++
	}
	fmt.Println("[utag] skipped:", skip)
	fmt.Println("[utag] inserted:", n)
}
