package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"call-for-white-papers/pkg/operations"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog"
)

var domains []string
var urls []*url.URL

func init() {
	// Get input through flags
	dataFilePath := flag.String("input", "input.json", "json configuration input file for seeding web crawler")
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Read Input
	jsonFile, err := os.Open(*dataFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read/Parse Input
	byteValue, _ := ioutil.ReadAll(jsonFile)
	i := &operations.Input{}
	json.Unmarshal(byteValue, i)
	for _, fullUrl := range i.Urls {
		u, err := url.Parse(fullUrl)
		if err != nil {
			os.Exit(1)
		}
		domains = append(domains, u.Hostname())
		urls = append(urls, u)
	}
}

func main() {

	// Create collector base on input
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(false),
		colly.AllowedDomains(domains...),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	c.SetRequestTimeout(15 * time.Second)

	// Run operations
	opt := operations.New(c, urls)
	opt.Start()
}
