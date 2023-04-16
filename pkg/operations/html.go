package operations

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

// logger object for html operations
var htmlHandler = log.With().Str("handler", "html").Logger()

// registerHtmlHandler adds all operations to run against an HTML on visit.
// Function will be executed on every HTML element matched by the GoQuery Selector parameter and the return time of handler
func (o *operations) registerHtmlHandler() {
	o.c.OnHTML(auditTitle(o))
	o.c.OnHTML(auditAndQueueHref(o))
}

// auditTitle runs against current title object when a 'title' element is found
func auditTitle(o *operations) (string, func(e *colly.HTMLElement)) {
	return "title", func(e *colly.HTMLElement) {
		htmlHandler.Debug().
			Uint32("id", e.Request.ID).
			Str("title", e.Text).
			Msg("title data")

		match, err := regexp.MatchString("(?i)white papers", e.Text)
		if err != nil {
			htmlHandler.Error().
				Uint32("id", e.Request.ID).
				Msg("failed to configure regex")
			return
		}

		// check if we have seen URL before and already flagged it for followup
		if _, ok := o.FlagForManualVisit[e.Request.URL.String()]; !ok && match {
			reason := "title match"
			htmlHandler.Info().
				Uint32("id", e.Request.ID).
				Str("reason", reason).
				Str("source", e.Request.URL.String()).
				Msg("Found opportunity")
			o.FlagForManualVisit[e.Request.URL.String()] = reason
		}
	}
}

// auditAndQueueHref validates all URLs found and registers possible new urls to be visited when 'a[href]' element is found
func auditAndQueueHref(o *operations) (string, func(e *colly.HTMLElement)) {
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

		// Remove fragments and queries from URL for simplicity sake
		// schema://host/path
		cleanedURL := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)

		if validateURL(o, u, e.Request.ID) {
			htmlHandler.Debug().
				Uint32("id", e.Request.ID).
				Str("source", u.String()).
				Str("cleaned", cleanedURL).
				Msg("registering new url to visit")

			o.c.Visit(cleanedURL)
		}
	}
}

func validateURL(o *operations, url *url.URL, requestID uint32) (shouldVisit bool) {
	// Check if URL has host name. If it doesn't we most likely don't want to visit it or can't
	if url.Hostname() == "" {
		htmlHandler.Debug().
			Uint32("id", requestID).
			Str("url", url.String()).
			Msg("url has no host")
		return
	}

	// Check to make sure we haven't visited the full url before to stop inf looping
	if _, ok := o.visitMap[url.String()]; !ok {
		o.visitMap[url.String()] = fmt.Sprintf("%d", requestID)
	} else {
		htmlHandler.Debug().
			Uint32("id", requestID).
			Str("url", url.String()).
			Msg("url already seen")
		return
	}

	// looks good, check URL out
	shouldVisit = true
	return
}
