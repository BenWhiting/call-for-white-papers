package operations

import (
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

// logger object for request operations
var requestHandler = log.With().Str("handler", "request").Logger()

// registerRequestHandler adds all operations to run before visiting a HTML page
func (o *operations) registerRequestHandler() {
	// register what to do on requests
	o.c.OnRequest(visitingMessenger(o))
}

func visitingMessenger(o *operations) func(*colly.Request) {
	return func(r *colly.Request) {
		requestHandler.Debug().
			Uint32("id", r.ID).
			Str("url", r.URL.String()).
			Msg("visiting target")
	}
}
