# botpublisher
[![Go Report Card](https://goreportcard.com/badge/github.com/junkichi/botpublisher)](https://goreportcard.com/report/github.com/junkichi/botpublisher) [![](https://godoc.org/github.com/junkichi/botpublisher?status.svg)](http://godoc.org/github.com/junkichi/botpublisher) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

`botpublisher` is twitter bot to publish the update of the site and tweets of hash tags.

## Requirements
- [MongoDB](https://www.mongodb.com/)
- [Google Chrome](https://www.google.co.jp/chrome/)

## Features
- collect the update of the site by RSS
- collect tweets of hash tags
- publish tweet of collections
- take screen shots of the site

## Dependencies
- [gofeed](https://github.com/mmcdole/gofeed) Parse RSS and Atom feeds in Go
- [Anaconda](https://github.com/ChimeraCoder/anaconda) A Go client library for the Twitter 1.1 API
- [chromedp](https://github.com/chromedp/chromedp) A faster, simpler way to drive browsers supporting the Chrome DevTools Protocol
- [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) The Go driver for MongoDB

## License
- This project is licensed under the [MIT License](https://github.com/junkichi/botpublisher/blob/master/LICENSE)
