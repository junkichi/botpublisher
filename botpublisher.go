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

	"botpublisher/storage"
	"botpublisher/twitter"
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

var publishCOL = "publish"

func publishDescription() {
	s := storage.GetInstance()
	desc, err := storage.FindPublish(s, publishCOL)
	if err != nil {
		return
	}
	fmt.Println("publish:", desc)
	tw := twitter.GetInstance()
	twitter.PublishTweet(tw, desc)
}

func publishWorker(publisherConfig PublisherConfig, ticker *time.Ticker, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()

	for {
		select {
		case <-stopCh:
			fmt.Println("publishWorker: stop request received.")
			return
		case t := <-ticker.C:
			fmt.Println("=== Publish <", t, "> ===")
			publishDescription()
		}
	}
}

type RssCollector interface {
	Init() bool
	Collect() bool
}

func collectWorker(publisherConfig PublisherConfig, ticker *time.Ticker, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()

	rssCollectors := []RssCollector{GoogleNewsRSS{}, UrayasuRSS{}}
	for _, rssCollector := range rssCollectors {
		result := rssCollector.Init()
		if result == false {
			return
		}
	}
	initUrayasuTagTweet(publisherConfig.UrayasuTagConfig.Query)

	for {
		select {
		case <-stopCh:
			fmt.Println("collectWorker: stop request received.")
			return
		case t := <-ticker.C:
			fmt.Println("=== Collect <", t, "> ===")
			for _, rssCollector := range rssCollectors {
				result := rssCollector.Collect()
				if result == false {
					fmt.Println("collectWorker: collect error")
				}
			}
			collectUrayasuTagTweet(publisherConfig.UrayasuTagConfig.Query)
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
	storage.SetConfig("mongodb://localhost:27017", "test")
	twitter.SetConfig(
		publisherConfig.TwitterConfig.AccessToken,
		publisherConfig.TwitterConfig.AccessTokenSecret,
		publisherConfig.TwitterConfig.ConsumerKey,
		publisherConfig.TwitterConfig.ConsumerSecret)

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
	storage.TermInstance()
}
