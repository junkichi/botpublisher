package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublisherConfig struct {
	IntervalConfig   `json:"interval"`
	TwitterConfig    `json:"twitter"`
	UrayasuTagConfig `json:"urayasuTag"`
}

type IntervalConfig struct {
	Collect int `json:"collect"`
	Publish int `json:"publish"`
}

type TwitterConfig struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
}

type UrayasuTagConfig struct {
	Query string `json:"query"`
}

var publishDbID = "test"
var publishDbCOL = "publish"

func publishDescription(c mongo.Client, api *anaconda.TwitterApi) {
	col := c.Database(publishDbID).Collection(publishDbCOL)
	desc, err := DBFindPublish(col)
	if err != nil {
		return
	}
	fmt.Println("publish:", desc)
	TWPublishTweet(api, desc)
}

func publishWorker(publisherConfig PublisherConfig, ticker *time.Ticker, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()

	api, _ := TWConnect(
		publisherConfig.TwitterConfig.AccessToken,
		publisherConfig.TwitterConfig.AccessTokenSecret,
		publisherConfig.TwitterConfig.ConsumerKey,
		publisherConfig.TwitterConfig.ConsumerSecret)

	c := DBConnect("mongodb://localhost:27017")

	for {
		select {
		case <-stopCh:
			fmt.Println("publishWorker: stop request received.")
			DBDisconnect(c)
			return
		case t := <-ticker.C:
			fmt.Println("=== Publish <", t, "> ===")
			publishDescription(c, api)
		}
	}
}

func collectWorker(publisherConfig PublisherConfig, ticker *time.Ticker, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()

	api, _ := TWConnect(
		publisherConfig.TwitterConfig.AccessToken,
		publisherConfig.TwitterConfig.AccessTokenSecret,
		publisherConfig.TwitterConfig.ConsumerKey,
		publisherConfig.TwitterConfig.ConsumerSecret)

	c := DBConnect("mongodb://localhost:27017")

	result := initUrayasuRSS(c)
	if result == false {
		DBDisconnect(c)
		return
	}
	result = initGoogleNewsRSS(c)
	if result == false {
		DBDisconnect(c)
		return
	}
	initUrayasuTagTweet(c, api, publisherConfig.UrayasuTagConfig.Query)

	for {
		select {
		case <-stopCh:
			fmt.Println("collectWorker: stop request received.")
			DBDisconnect(c)
			return
		case t := <-ticker.C:
			fmt.Println("=== Collect <", t, "> ===")
			result := collectUrayasuRSS(c)
			if result == false {
				DBDisconnect(c)
				return
			}
			result = collectGoogleNewsRSS(c)
			if result == false {
				DBDisconnect(c)
				return
			}
			collectUrayasuTagTweet(c, api, publisherConfig.UrayasuTagConfig.Query)
		}
	}
}

func setupCloseHandler(doneCh chan struct{}) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		close(doneCh)
	}()
}

func main() {
	confFile, err := ioutil.ReadFile("botpublisher.json")
	if err != nil {
		fmt.Println("read error:", err)
		return
	}

	var publisherConfig PublisherConfig
	json.Unmarshal(confFile, &publisherConfig)

	doneCh := make(chan struct{})
	stopCh := make(chan struct{})
	wg := sync.WaitGroup{}

	collectTicker := time.NewTicker(time.Duration(publisherConfig.IntervalConfig.Collect) * time.Second)
	wg.Add(1)
	go collectWorker(publisherConfig, collectTicker, stopCh, &wg)

	publishTicker := time.NewTicker(time.Duration(publisherConfig.IntervalConfig.Publish) * time.Second)
	wg.Add(1)
	go publishWorker(publisherConfig, publishTicker, stopCh, &wg)

	setupCloseHandler(doneCh)

	<-doneCh
	collectTicker.Stop()
	publishTicker.Stop()
	close(stopCh)
	wg.Wait()
}
