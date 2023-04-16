package operations

import (
	"net/url"

	"github.com/gocolly/colly"
)

type Input struct {
	Urls []string `json:"urls"`
}

type operations struct {
	c                  *colly.Collector
	urls               []*url.URL
	visitMap           map[string]string
	FlagForManualVisit map[string]string
}

func New(c *colly.Collector, u []*url.URL) *operations {
	return &operations{
		c:                  c,
		urls:               u,
		visitMap:           make(map[string]string),
		FlagForManualVisit: make(map[string]string),
	}
}

func (o *operations) Start() error {
	o.registerHtmlHandler()
	o.registerRequestHandler()
	o.registerResponseHandler()

	// Loop over every base URL we are given at start time and visit it.
	for _, url := range o.urls {
		o.c.Visit(url.String())
	}

	o.c.Wait()
	return nil
}
