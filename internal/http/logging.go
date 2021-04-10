package http

import (
	"net/http"
	"net/http/httputil"
	"time"

	log "github.com/sirupsen/logrus"
)

// NewLoggingRoundTripper creates a new logging roundtripper
func NewLoggingRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return LoggingRoundTripper{
		Next: rt,
	}
}

// LoggingRoundTripper logs a request-response
type LoggingRoundTripper struct {
	Next http.RoundTripper
}

// RoundTrip logs a request-response
func (l LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	req, _ := httputil.DumpRequestOut(r, true)

	var roundTripper http.RoundTripper
	if l.Next != nil {
		roundTripper = l.Next
	} else {
		roundTripper = http.DefaultTransport
	}

	start := time.Now()
	resp, err := roundTripper.RoundTrip(r)
	took := time.Since(start)

	var res []byte
	if resp != nil {
		res, _ = httputil.DumpResponse(resp, true)
	}

	logger := log.WithFields(log.Fields{
		"host":     r.Host,
		"took":     took,
		"request":  string(req),
		"response": string(res),
	})
	logger.Trace("http request")

	return resp, err
}
