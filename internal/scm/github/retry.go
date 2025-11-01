package github

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v76/github"
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

type retryAfterError time.Duration

func (r retryAfterError) Error() string {
	return fmt.Sprintf("rate limit exceeded, waiting for %s)", time.Duration(r))
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
	var response *github.Response
	err := retryAPIRequest(ctx, func() error {
		var err error
		response, err = fn()

		if response != nil {
			httpResponse := response.Response
			retryAfterErr := retryAfterFromHTTPResponse(httpResponse)
			if retryAfterErr != nil {
				return retryAfterErr
			}
		}

		return err
	})
	return response, err
}

func retryAPIRequest(ctx context.Context, fn func() error) error {
	tries := 0

	for {
		tries++

		err := fn()
		if err == nil { // NB!
			return nil
		}

		var retryAfter retryAfterError
		switch {
		// If GitHub has specified how long we should wait, use that information
		case errors.As(err, &retryAfter):
			err := sleep(ctx, time.Duration(retryAfter))
			if err != nil {
				return err
			}
		// If secondary rate limit error, use an exponential back-off to determine the wait
		case strings.Contains(err.Error(), "secondary rate limit"):
			err := sleep(ctx, exponentialBackoff(tries))
			if err != nil {
				return err
			}
		// If any other error, return the error
		default:
			return err
		}
	}
}

func retryAfterFromHTTPResponse(response *http.Response) error {
	if response == nil {
		return nil
	}

	retryAfterStr := response.Header.Get(retryHeader)
	if retryAfterStr == "" {
		return nil
	}

	retryAfterSeconds, err := strconv.Atoi(retryAfterStr)
	if err != nil {
		return nil
	}

	if retryAfterSeconds <= 0 {
		return nil
	}

	return retryAfterError(time.Duration(retryAfterSeconds) * time.Second)
}

func exponentialBackoff(tries int) time.Duration {
	// 10, 80, 270... seconds
	return time.Duration(math.Pow(float64(tries), 3)) * 10 * time.Second
}
