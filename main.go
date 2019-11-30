package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/gocolly/colly"

	"github.com/sirupsen/logrus"
)

var TO *url.URL

type Client struct {
	URL   *url.URL
	Colly *colly.Collector
}

func New(to *url.URL, c *colly.Collector) *Client {
	return &Client{
		URL:   to,
		Colly: c,
	}
}

func init() {
	to := flag.String("to", "https://google.com", "gets [to's] documents")
	flag.Parse()

	u, err := url.Parse(*to)
	if err != nil {
		logrus.Fatal(err)
	}
	TO = u
}

func main() {
	c := colly.NewCollector()
	client := New(TO, c)

	f, err := os.OpenFile("./raw.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()

	// Find and visit all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		html := string(e.Response.Body)
		f.Write([]byte(html))
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	client.Colly.Visit(client.URL.String())

}
