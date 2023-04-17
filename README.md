# call-for-white-papers
golang based web crawler scraping websites for possible call for white papers.

## What is a Call for White paper?

A white paper is an informational document issued by a company or not-for-profit organization 
to promote or highlight the features of a solution, product, or service that it offers or plans to offer.

When an organization puts out a "call for white papers" it means the organization is looking for third 
parties to provide a white paper to solve one of their needs.

### What is this project solving?

Finding possible call for white papers can be a timely process. This repository is designed to provide 
a cli & library to automatically scrape websites for calls for white papers.

## Getting started
 
 This repository provides a VScode developer container so all is needed is [VScode](https://code.visualstudio.com/) and [Docker](https://www.docker.com/) enabled.

### Build

Navigate to [cmd](/cmd/) and there you will see a makefile. 
Run the command `make build` to create the cli version of the codebase

### Run

Once the code is compiled run the code like so:

``` bash
./cmd/build/papers-please --input ./data/urls_single.json -debug --max-depth 7 --concurrent 2
```

### Input

As seen under [data](/data/) the cli takes in a json formatted file for all base urls to scan and start crawling from.

The code expects a json file of the following schema:

``` json
{
    "urls": [      
        // list of urls
    ]
}
```