package github

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const retryHeader = "Retry-After"

var sleep = func(ctx context.Context, d time.Duration) {
	log.Infof("Hit rate limit, sleeping for %s", d)
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}

// retry runs a GitHub API request and retries it if a temporary error occurred
func retry[K any](ctx context.Context, fn func() (K, *github.Response, error)) (K, *github.Response, error) {
	var val K
	resp, err := retryWithoutReturn(ctx, func() (*github.Response, error) {
		var resp *github.Response
		var err error
		val, resp, err = fn()
		return resp, err
	})
	return val, resp, err
}

// retry runs a GitHub API request with no return value and retries it if a temporary error occurred
func retryWithoutReturn(ctx context.Context, fn func() (*github.Response, error)) (*github.Response, error) {
	tries := 0

	for {
		tries++

		httpResp, err := fn()
		if err == nil { // NB!
			return httpResp, nil
		}

		retryAfterStr := httpResp.Header.Get(retryHeader)

		if retryAfterStr == "" && !strings.Contains(err.Error(), "secondary rate limit") {
			return httpResp, err
		}

		// If GitHub has specified how long we should wait, use that information
		if httpResp != nil && httpResp.Header != nil {
			if retryAfterStr != "" {
				if retryAfter, err := strconv.Atoi(retryAfterStr); err != nil {
					sleep(ctx, time.Duration(retryAfter)*time.Second)
					continue
				} else {
					return httpResp, errors.WithMessage(err, "could not convert Retry-After header")
				}
			}
		}

		// Otherwise use an exponential back-off
		sleep(ctx, time.Duration(math.Pow(float64(tries), 3))*10*time.Second)
	}
}
