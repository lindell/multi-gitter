package github

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const retryHeader = "Retry-After"

var sleep = func(ctx context.Context, d time.Duration) error {
	log.Infof("Hit rate limit, sleeping for %s", d)
	select {
	case <-ctx.Done():
		return errors.New("aborted while waiting for rate-limit")
	case <-time.After(d):
		return nil
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

// retryWithoutReturn runs a GitHub API request with no return value and retries it if a temporary error occurred
func retryWithoutReturn(ctx context.Context, fn func() (*github.Response, error)) (*github.Response, error) {
	tries := 0

	for {
		tries++

		githubResp, err := fn()
		if err == nil { // NB!
			return githubResp, nil
		}

		// Get the number of retry seconds (if any)
		retryAfter := 0
		if githubResp != nil && githubResp.Header != nil {
			retryAfterStr := githubResp.Header.Get(retryHeader)
			if retryAfterStr != "" {
				var err error
				if retryAfter, err = strconv.Atoi(retryAfterStr); err != nil {
					return githubResp, errors.WithMessage(err, "could not convert Retry-After header")
				}
			}
		}

		switch {
		// If GitHub has specified how long we should wait, use that information
		case retryAfter != 0:
			err := sleep(ctx, time.Duration(retryAfter)*time.Second)
			if err != nil {
				return githubResp, err
			}
		// If secondary rate limit error, use an exponential back-off to determine the wait
		case strings.Contains(err.Error(), "secondary rate limit"):
			err := sleep(ctx, time.Duration(math.Pow(float64(tries), 3))*10*time.Second)
			if err != nil {
				return githubResp, err
			}
		// If any other error, return the error
		default:
			return githubResp, err
		}
	}
}
