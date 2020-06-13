package rss

import (
	"errors"
	"fmt"

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
