package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"papers-please/pkg/operations"

	"github.com/gocolly/colly"
)

func main() {
	// Get Input
	jsonFile, err := os.Open("../data/urls.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read/Parse Input
	byteValue, _ := ioutil.ReadAll(jsonFile)
	i := &operations.Input{}
	json.Unmarshal(byteValue, i)
	d := []string{}
	us := []*url.URL{}
	for _, fullUrl := range i.Urls {
		u, err := url.Parse(fullUrl)
		if err != nil {
			os.Exit(1)
		}
		d = append(d, u.Hostname())
		us = append(us, u)
	}

	// Create collector base on input
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(false),
		colly.AllowedDomains(d...),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// Do the thing
	opt := operations.New(c, us)
	opt.Start()
	time.Sleep(8 * time.Second)
}
