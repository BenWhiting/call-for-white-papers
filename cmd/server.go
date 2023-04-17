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
var concurrent *int
var maxDepth *int

func init() {
	dataFilePath := flag.String("input", "input.json", "json configuration input file for seeding web crawler")
	debug := flag.Bool("debug", false, "sets log level to debug")
	concurrent = flag.Int("concurrent", 2, "max number of concurrent processes while running")
	maxDepth = flag.Int("max-depth", 1, "max-depth limits the recursion depth of visited URLs.")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Read Input
	jsonFile, err := os.Open(*dataFilePath)
	if err != nil {
		panic(err)
	}

	// Read/Parse Input
	byteValue, _ := ioutil.ReadAll(jsonFile)
	i := &operations.Input{}
	json.Unmarshal(byteValue, i)
	for _, fullUrl := range i.Urls {
		u, err := url.Parse(fullUrl)
		if err != nil {
			panic(err)
		}
		domains = append(domains, u.Hostname())
		urls = append(urls, u)
	}
}

func main() {
	// Create collector base on input
	c := colly.NewCollector(
		colly.MaxDepth(*maxDepth),
		colly.Async(false),
		colly.AllowedDomains(domains...),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: *concurrent})
	c.SetRequestTimeout(15 * time.Second)

	// Run operations
	opt := operations.New(c, urls)
	opt.Start()
	fmt.Println("I-------")
	for key, value := range opt.FlagForManualVisit {
		fmt.Printf("Checkout out %s. Reason: %s\n", key, value)
	}

}
