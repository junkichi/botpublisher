package main

import (
	"errors"
	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

func TWConnect(accessToken string, accessTokenSecret string,
	consumerKey string, consumerSecret string) (*anaconda.TwitterApi, error) {
	api := anaconda.NewTwitterApiWithCredentials(
		accessToken,
		accessTokenSecret,
		consumerKey,
		consumerSecret)
	if api == nil {
		fmt.Println("newapi with credentials error")
		return nil, errors.New("credentials error")
	}

	return api, nil
}

func TWSearchTweet(api *anaconda.TwitterApi, query string, tweets *[]anaconda.Tweet) {
	searchResult, _ := api.GetSearch(query, nil)
	*tweets = searchResult.Statuses
}

func TWPublishTweet(api *anaconda.TwitterApi, desc string) {
	_, err := api.PostTweet(desc, nil)
	if err != nil {
		fmt.Println("post tweet error:", err)
	}
}
