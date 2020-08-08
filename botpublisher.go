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

	"github.com/junkichi/botpublisher/storage"
	"github.com/junkichi/botpublisher/twitter"
)

// PublisherConfig is the configrations of botpublisher
type PublisherConfig struct {
	IntervalConfig   `json:"interval"`
	TwitterConfig    `json:"twitter"`
	StorageConfig    `json:"storage"`
	UrayasuTagConfig `json:"urayasuTag"`
	ScreenShotConfig `json:"screenShot"`
}

// IntervalConfig is the configration of processing interval
type IntervalConfig struct {
	Collect int `json:"collect"`
	Publish int `json:"publish"`
}

// TwitterConfig is the configration of tweet api
type TwitterConfig struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
}

// StorageConfig is the configration of database
type StorageConfig struct {
	Url    string `json:"url"`
	Id     string `json:"id"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

// UrayasuTagConfig is the configration of twitter search
type UrayasuTagConfig struct {
	Query string `json:"query"`
}

// ScreenShotConfig is the configration of image directory
type ScreenShotConfig struct {
	ImageDir string `json:"imagedir"`
}

var publishCOL = "publish"

func publishDescription() {
	s := storage.GetInstance()
	desc, imgpath, err := storage.FindPublish(s, publishCOL)
	if err != nil {
		return
	}
	fmt.Println("publish:", desc, imgpath)
	tw := twitter.GetInstance()
	twitter.PublishTweet(tw, desc, imgpath)
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

// RssCollector is the interface of collecting RSS
type RssCollector interface {
	GetImageDir() string
	Init() bool
	Collect(string) bool
}

func rssCollectWorker(rssCollectors []RssCollector, ticker *time.Ticker, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()

	for _, rssCollector := range rssCollectors {
		result := rssCollector.Init()
		if result == false {
			return
		}
	}

	for {
		select {
		case <-stopCh:
			fmt.Println("rssCollectWorker: stop request received.")
			return
		case t := <-ticker.C:
			fmt.Println("=== RSS Collect <", t, "> ===")
			for _, rssCollector := range rssCollectors {
				result := rssCollector.Collect(rssCollector.GetImageDir())
				if result == false {
					fmt.Println("collectWorker: collect error")
				}
			}
		}
	}
}

// TweetCollector is the interface of collecting Tweet
type TweetCollector interface {
	GetQuery() string
	Init(string)
	Collect(string)
}

func tweetCollectWorker(twCollectors []TweetCollector, ticker *time.Ticker, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()

	for _, twCollector := range twCollectors {
		twCollector.Init(twCollector.GetQuery())
	}

	for {
		select {
		case <-stopCh:
			fmt.Println("tweetCollectWorker: stop request received.")
			return
		case t := <-ticker.C:
			fmt.Println("=== Tweet Collect <", t, "> ===")
			for _, twCollector := range twCollectors {
				twCollector.Collect(twCollector.GetQuery())
			}
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
	storage.SetConfig(
		publisherConfig.StorageConfig.Url,
		publisherConfig.StorageConfig.Id,
		publisherConfig.StorageConfig.User,
		publisherConfig.StorageConfig.Passwd)
	twitter.SetConfig(
		publisherConfig.TwitterConfig.AccessToken,
		publisherConfig.TwitterConfig.AccessTokenSecret,
		publisherConfig.TwitterConfig.ConsumerKey,
		publisherConfig.TwitterConfig.ConsumerSecret)

	doneCh := make(chan struct{})
	stopCh := make(chan struct{})
	wg := sync.WaitGroup{}

	rssCollectors := []RssCollector{GoogleNewsRSS{""}, UrayasuRSS{publisherConfig.ScreenShotConfig.ImageDir}}
	rssCollectTicker := time.NewTicker(time.Duration(publisherConfig.IntervalConfig.Collect) * time.Second)
	wg.Add(1)
	go rssCollectWorker(rssCollectors, rssCollectTicker, stopCh, &wg)

	twCollectors := []TweetCollector{UrayasuTagTweet{publisherConfig.UrayasuTagConfig.Query}}
	twCollectTicker := time.NewTicker(time.Duration(publisherConfig.IntervalConfig.Collect) * time.Second)
	wg.Add(1)
	go tweetCollectWorker(twCollectors, twCollectTicker, stopCh, &wg)

	publishTicker := time.NewTicker(time.Duration(publisherConfig.IntervalConfig.Publish) * time.Second)
	wg.Add(1)
	go publishWorker(publisherConfig, publishTicker, stopCh, &wg)

	setupCloseHandler(doneCh)

	<-doneCh
	rssCollectTicker.Stop()
	twCollectTicker.Stop()
	publishTicker.Stop()
	close(stopCh)
	wg.Wait()
	storage.TermInstance()
}
