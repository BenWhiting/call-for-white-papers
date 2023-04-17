package operations

import (
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

// logger object for response operations
var responseHandler = log.With().Str("handler", "response").Logger()

// registerResponseHandler adds all operations to run after after running a HTTP GET command against a URL
func (o *operations) registerResponseHandler() {
	o.c.OnResponse(printStatusCode(o))
	o.c.OnError(responseErrorHandler(o))
}

// printStatusCode prints the returned status code
func printStatusCode(o *operations) func(*colly.Response) {
	return func(r *colly.Response) {
		responseHandler.Debug().
			Uint32("id", r.Request.ID).
			Int("code", r.StatusCode).
			Str("target", r.Request.URL.String()).
			Msg("target response")
	}
}

// responseErrorHandler handles when a request returns an error
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
