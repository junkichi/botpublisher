package rss

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/mmcdole/gofeed"
)

func RetrieveRSS(rssurl string, items *[]*gofeed.Item) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssurl)
	if err != nil {
		fmt.Println("fp.ParseURL() Error:", err)
		return errors.New("dont parse url. " + rssurl)
	}
	*items = feed.Items
	return nil
}

func TakeScreenShot(url, sel, imgpath string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx, elementScreenshot(url, sel, &buf))
	if err != nil {
		fmt.Println("chrome run error:", err)
		return err
	}

	err = ioutil.WriteFile(imgpath, buf, 0644)
	if err != nil {
		fmt.Println("write file error:", err)
		return err
	}

	return nil
}

func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	width, height := 1920, 10000
	return chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false),
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByID),
	}
}
