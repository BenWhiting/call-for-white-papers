package operations

import (
	"net/url"
	"regexp"
	"strings"

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

		// check if we have seen URL before and already flagged it for followup
		if _, ok := o.FlagForManualVisit[e.Request.URL.String()]; !ok && callsForPaper(e) {
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

// callsForPaper checks the title object for possible call for white paper text
func callsForPaper(e *colly.HTMLElement) bool {
	whitePaper, err := regexp.MatchString("(?i)white papers", e.Text)
	if err != nil {
		htmlHandler.Error().
			Uint32("id", e.Request.ID).
			Msg("failed to configure regex")
		return false
	}

	paper, err := regexp.MatchString("(?i)call for papers", e.Text)
	if err != nil {
		htmlHandler.Error().
			Uint32("id", e.Request.ID).
			Msg("failed to configure regex")
		return false
	}

	return whitePaper || paper
}

// auditAndQueueHref validates all URLs found and registers possible new urls to be visited when 'a[href]' element is found
func auditAndQueueHref(o *operations) (string, func(e *colly.HTMLElement)) {
	return "a[href]", func(e *colly.HTMLElement) {
		urlString := e.Request.AbsoluteURL(e.Attr("href"))
		if validateURLForVisit(urlString, o) {
			htmlHandler.Debug().
				Uint32("id", e.Request.ID).
				Str("source", urlString).
				Msg("registering new url to visit")
			// we have now seen
			o.hasRegistered[urlString] = true
			// prep for visit
			err := e.Request.Visit(urlString)
			if err != nil {
				htmlHandler.Error().
					Uint32("id", e.Request.ID).
					Str("source", urlString).
					Err(err)
				return
			}
		} else {
			htmlHandler.Debug().
				Uint32("id", e.Request.ID).
				Str("source", urlString).
				Msg("don't need to visit.")
		}
	}
}

// validate URL before putting into the queue
func validateURLForVisit(u string, o *operations) (shouldVisit bool) {
	// empty url isn't valid
	if u == "" {
		return
	}

	// Check if we have already seen it and registered it.
	// Even if it is not a valid url to search, register to skip later
	if _, seen := o.hasRegistered[u]; seen {
		return
	}

	// get url object
	parsed, err := url.Parse(u)
	if err != nil {
		return
	}

	// Don't visit URLs with fragments
	if len(parsed.Fragment) > 1 {
		return
	}

	// Call for white papers are never under the following segments
	var segmentIgnores = [...]string{"uploads", "documents"}
	parsedSplit := strings.Split(parsed.Path, "/")
	for _, segment := range parsedSplit {
		for _, ignores := range segmentIgnores {
			if segment == ignores {
				return
			}
		}
	}

	shouldVisit = true
	return
}
