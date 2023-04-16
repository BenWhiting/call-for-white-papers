package operations

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

var htmlHandler = log.With().Str("handler", "html").Logger()

func (o *operations) registerHtmlHandler() {
	// register what to do on ever HTML element with the selector
	// Print out title
	o.c.OnHTML(printTitle(o))
	o.c.OnHTML(queueHref(o))
}

func printTitle(o *operations) (string, func(e *colly.HTMLElement)) {
	return "title", func(e *colly.HTMLElement) {
		if _, ok := o.visitMap[e.Request.URL.String()]; ok {
			// do something here
		}
		htmlHandler.Debug().
			Uint32("id", e.Request.ID).
			Str("title", e.Text).
			Msg("title data")
	}
}

func queueHref(o *operations) (string, func(e *colly.HTMLElement)) {
	return "a[href]", func(e *colly.HTMLElement) {
		urlString := e.Attr("href")

		// Parse to URL type
		u, err := url.Parse(urlString)
		if err != nil {
			htmlHandler.Error().
				Uint32("id", e.Request.ID).
				Str("url", u.String()).
				Msg("new url not valid")
			return
		}
		// ##### Validation of URL #####
		// check for  host name
		if u.Hostname() == "" {
			htmlHandler.Debug().
				Uint32("id", e.Request.ID).
				Str("url", u.String()).
				Msg("url has no host")
			return
		}
		// make sure we don't double dip
		if _, ok := o.visitMap[u.String()]; !ok {
			o.visitMap[u.String()] = fmt.Sprintf("%d", e.Request.ID)
		} else {
			htmlHandler.Debug().
				Uint32("id", e.Request.ID).
				Str("url", u.String()).
				Msg("url already seen")
			return
		}
		// ##### End Validation #####
		htmlHandler.Debug().
			Uint32("id", e.Request.ID).
			Str("url", u.String()).
			Msg("registering new url to visit")

		o.c.Visit(u.String())
	}
}
