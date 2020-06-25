package github

import (
	"net/http"
	"net/http/httputil"
	"time"

	log "github.com/sirupsen/logrus"
)

type loggingRoundTripper struct {
	next http.RoundTripper
}

func (l loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	req, _ := httputil.DumpRequestOut(r, true)

	start := time.Now()
	resp, err := l.next.RoundTrip(r)
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
