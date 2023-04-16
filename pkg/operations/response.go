package operations

import (
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

var responseHandler = log.With().Str("handler", "response").Logger()

func (o *operations) registerResponseHandler() {
	// regester what to do on response
	o.c.OnResponse(printStatusCode(o))
	o.c.OnError(responseErrorHandler(o))
}

func printStatusCode(o *operations) func(*colly.Response) {
	return func(r *colly.Response) {
		responseHandler.Debug().
			Uint32("id", r.Request.ID).
			Int("code", r.StatusCode).
			Msg("target response")
	}
}

// Set error handler
func responseErrorHandler(o *operations) func(*colly.Response, error) {
	return func(r *colly.Response, err error) {
		responseHandler.Error().
			Uint32("id", r.Request.ID).
			Int("code", r.StatusCode).
			Str("url", r.Request.URL.String()).
			Err(err).
			Msg("failed to get response")
	}
}
