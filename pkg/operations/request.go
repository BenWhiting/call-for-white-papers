package operations

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

// logger object for request operations
var requestHandler = log.With().Str("handler", "request").Logger()

// registerRequestHandler adds all operations to run before visiting a HTML page
func (o *operations) registerRequestHandler() {
	// register what to do on requests
	o.c.OnRequest(hasVisitedCircuitBreaker(o))
}

// Before running requests check if we need to run the request
func hasVisitedCircuitBreaker(o *operations) func(*colly.Request) {
	return func(r *colly.Request) {
		// Check to make sure we haven't visited/processed the full url before to stop inf looping
		if _, ok := o.hasVisited[r.URL.String()]; !ok {
			o.hasVisited[r.URL.String()] = fmt.Sprintf("%d", r.ID)
		} else {
			htmlHandler.Debug().
				Uint32("id", r.ID).
				Str("url", r.URL.String()).
				Msg("url already visited")
			r.Abort()
			return
		}

		requestHandler.Debug().
			Uint32("id", r.ID).
			Str("url", r.URL.String()).
			Msg("visiting URL")
	}
}
