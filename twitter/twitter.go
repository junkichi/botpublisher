package twitter

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"

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

func PublishTweet(tw *Twitter, desc string, imgpath string) {
	v := url.Values{}
	if len(imgpath) == 0 {
		v = nil
	} else {
		base64String, err := convertImage2Base64(imgpath)
		if err != nil {
			v = nil
		} else {
			media, _ := tw.api.UploadMedia(base64String)
			v.Add("media_ids", media.MediaIDString)
		}
	}
	_, err := tw.api.PostTweet(desc, v)
	if err != nil {
		fmt.Println("post tweet error:", err)
	}
}

func convertImage2Base64(imgpath string) (string, error) {
	data, err := ioutil.ReadFile(imgpath)
	if err != nil {
		fmt.Println("read image error:", err)
		return "", err
	}
	base64String := base64.StdEncoding.EncodeToString(data)
	return base64String, nil
}
