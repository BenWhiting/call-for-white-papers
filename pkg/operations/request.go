package operations

import (
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

var requestHandler = log.With().Str("handler", "request").Logger()

func (o *operations) registerRequestHandler() {
	// regester what to do on requests
	o.c.OnRequest(VisitingMessager(o))
}

func VisitingMessager(o *operations) func(*colly.Request) {
	return func(r *colly.Request) {
		requestHandler.Debug().
			Uint32("id", r.ID).
			Str("url", r.URL.String()).
			Msg("visiting target")
	}
}
