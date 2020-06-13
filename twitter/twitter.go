package twitter

import (
	"errors"
	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

type Twitter struct {
	api *anaconda.TwitterApi
}

var twAccessToken string
var twAccessTokenSecret string
var twConsumerKey string
var twConsumerSecret string
var sharedInstance *Twitter

func newTwitter() *Twitter {
	if twAccessToken == "" {
		return nil
	}

	fmt.Println("newTwitter: success")
	api, _ := connect(twAccessToken, twAccessTokenSecret, twConsumerKey, twConsumerSecret)
	return &Twitter{api}
}

func SetConfig(accessToken string, accessTokenSecret string, consumerKey string, consumerSecret string) {
	twAccessToken = accessToken
	twAccessTokenSecret = accessTokenSecret
	twConsumerKey = consumerKey
	twConsumerSecret = consumerSecret
}

func GetInstance() *Twitter {
	if sharedInstance == nil {
		sharedInstance = newTwitter()
	}
	return sharedInstance
}

func connect(accessToken string, accessTokenSecret string,
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

func SearchTweet(tw *Twitter, query string, tweets *[]anaconda.Tweet) {
	searchResult, _ := tw.api.GetSearch(query, nil)
	*tweets = searchResult.Statuses
}

func PublishTweet(tw *Twitter, desc string) {
	_, err := tw.api.PostTweet(desc, nil)
	if err != nil {
		fmt.Println("post tweet error:", err)
	}
}
