package github

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_retryWithoutReturn(t *testing.T) {
	tests := []struct {
		name      string
		responses []response
		wantErr   bool
		sleep     time.Duration
	}{
		{
			name: "one fail",
			responses: []response{
				secondaryRateLimitError,
				okResponse,
			},
			sleep: 10 * time.Second,
		},
		{
			name: "two fails",
			responses: []response{
				secondaryRateLimitError,
				secondaryRateLimitError,
				okResponse,
			},
			sleep: 1*time.Minute + 30*time.Second,
		},
		{
			name: "three fails",
			responses: []response{
				secondaryRateLimitError,
				secondaryRateLimitError,
				secondaryRateLimitError,
				okResponse,
			},
			sleep: 6*time.Minute + 0*time.Second,
		},
		{
			name: "four fails",
			responses: []response{
				secondaryRateLimitError,
				secondaryRateLimitError,
				secondaryRateLimitError,
				secondaryRateLimitError,
				okResponse,
			},
			sleep: 16*time.Minute + 40*time.Second,
		},
		{
			name: "a real fail",
			responses: []response{
				realError,
			},
			sleep:   0,
			wantErr: true,
		},
		{
			name: "two temporary fails and a real one",
			responses: []response{
				secondaryRateLimitError,
				secondaryRateLimitError,
				realError,
			},
			sleep:   1*time.Minute + 30*time.Second,
			wantErr: true,
		},
		{
			name: "primary rate limit",
			responses: []response{
				primaryRateLimitError,
				okResponse,
			},
			sleep:   1*time.Minute + 37*time.Second,
			wantErr: false,
		},
		{
			name: "two primary rate limit errors in a row",
			responses: []response{
				primaryRateLimitError,
				primaryRateLimitError,
				okResponse,
			},
			sleep:   3*time.Minute + 14*time.Second,
			wantErr: false,
		},
		{
			name: "secondary rate limit followed by primary rate limit",
			responses: []response{
				secondaryRateLimitError,
				primaryRateLimitError,
				okResponse,
			},
			sleep:   1*time.Minute + 47*time.Second,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var slept time.Duration
			sleep = func(ctx context.Context, d time.Duration) error {
				slept += d
				return nil
			}

			call := 0
			fn := func() (*github.Response, error) {
				resp := tt.responses[call]
				call++
				return resp.response, resp.err
			}

			if _, err := retryWithoutReturn(context.Background(), fn); (err != nil) != tt.wantErr {
				t.Errorf("retry() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.sleep, slept)
		})
	}
}

func Test_retry(t *testing.T) {
	var slept time.Duration
	sleep = func(ctx context.Context, d time.Duration) error {
		slept += d
		return nil
	}

	call := 0
	fn := func() (*github.PullRequest, *github.Response, error) {
		call++
		if call == 4 {
			return &github.PullRequest{
				ID: &[]int64{100}[0],
			}, okResponse.response, nil
		}
		return nil, secondaryRateLimitError.response, secondaryRateLimitError.err
	}

	pr, resp, err := retry(context.Background(), fn)
	assert.Equal(t, int64(100), *pr.ID)
	assert.Equal(t, 200, resp.StatusCode)
	assert.NoError(t, err)
}

type response struct {
	response *github.Response
	err      error
}

var okResponse = response{
	err: nil,
	response: &github.Response{
		Response: &http.Response{
			StatusCode: 200,
		},
	},
}
var realError = createResponse(http.StatusNotFound, errors.New("something went wrong"))
var secondaryErrorMsg = "You have exceeded a secondary rate limit and have been temporarily blocked from content creation. Please retry your request again later."
var secondaryRateLimitError = createResponse(http.StatusForbidden, errors.New(secondaryErrorMsg))
var primaryRateLimitError = response{
	response: &github.Response{
		Response: &http.Response{
			StatusCode: http.StatusForbidden,
			Header: http.Header{
				"Retry-After": []string{"97"},
			},
		},
	},
	err: errors.New("rate limited"),
}

func createResponse(statusCode int, err error) response {
	return response{
		response: createGithubResponse(statusCode, err.Error()),
		err:      err,
	}
}

func createGithubResponse(statusCode int, errorMsg string) *github.Response {
	return &github.Response{
		Response: &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(strings.NewReader(errorMsg)),
		},
	}
}
