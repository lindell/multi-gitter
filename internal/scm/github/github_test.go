package github_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/lindell/multi-gitter/internal/scm/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTransport struct {
	pathBodies map[string]string
}

func (tt testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, ok := tt.pathBodies[req.URL.Path]
	if !ok {
		return nil, errors.New("could not find path")
	}
	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Request:       req,
		Header:        make(http.Header),
	}, nil
}

func (tt testTransport) Wrapper(http.RoundTripper) http.RoundTripper {
	return tt
}

func Test_GetRepositories(t *testing.T) {
	transport := testTransport{
		pathBodies: map[string]string{
			"/orgs/test-org/repos": `[
				{
					"id": 1,
					"name": "test1",
					"full_name": "test-org/test1",
					"private": false,
					"topics": [
						"frontend"
					],
					"owner": {
						"login": "test-org",
						"type": "Organization",
						"site_admin": false
					},
					"html_url": "https://github.com/test-org/test1",
					"fork": false,
					"archived": false,
					"disabled": false,
					"default_branch": "master",
					"permissions": {
						"admin": true,
						"push": true,
						"pull": true
					},
					"created_at": "2020-01-01T16:49:16Z"
				}
			]`,
			"/repos/test-org/test1": `{
				"id": 2,
				"name": "test1",
				"full_name": "test-org/test1",
				"private": false,
				"owner": {
					"login": "test-org",
					"type": "Organization",
					"site_admin": false
				},
				"html_url": "https://github.com/test-org/test1",
				"fork": false,
				"archived": false,
				"disabled": false,
				"default_branch": "master",
				"permissions": {
					"admin": true,
					"push": true,
					"pull": true
				},
				"created_at": "2020-01-02T16:49:16Z"
			}`,
			"/users/test-user/repos": `[
				{
					"id": 3,
					"name": "test2",
					"full_name": "lindell/test2",
					"private": false,
					"topics": [
						"backend",
						"go"
					],
					"owner": {
						"login": "lindell",
						"type": "User",
						"site_admin": false
					},
					"html_url": "https://github.com/lindell/test2",
					"fork": true,
					"archived": false,
					"disabled": false,
					"default_branch": "main",
					"permissions": {
						"admin": true,
						"push": true,
						"pull": true
					},
					"created_at": "2020-01-03T16:49:16Z"
				}
			]`,
		},
	}

	// Organization
	{
		gh, err := github.New(github.Config{
			TransportMiddleware: transport.Wrapper,
			RepoListing: github.RepositoryListing{
				Organizations: []string{"test-org"},
			},
			MergeTypes: []scm.MergeType{scm.MergeTypeMerge},
		})
		require.NoError(t, err)

		repos, err := gh.GetRepositories(context.Background())
		assert.NoError(t, err)
		if assert.Len(t, repos, 1) {
			assert.Equal(t, "master", repos[0].DefaultBranch())
			assert.Equal(t, "test-org/test1", repos[0].FullName())
		}
	}

	// repository
	{
		gh, err := github.New(github.Config{
			TransportMiddleware: transport.Wrapper,
			RepoListing: github.RepositoryListing{
				Repositories: []github.RepositoryReference{
					{
						OwnerName: "test-org",
						Name:      "test1",
					},
				},
			},
			MergeTypes: []scm.MergeType{scm.MergeTypeMerge},
		})
		require.NoError(t, err)

		repos, err := gh.GetRepositories(context.Background())
		assert.NoError(t, err)
		if assert.Len(t, repos, 1) {
			assert.Equal(t, "master", repos[0].DefaultBranch())
			assert.Equal(t, "test-org/test1", repos[0].FullName())
		}
	}

	// User
	{
		gh, err := github.New(github.Config{
			TransportMiddleware: transport.Wrapper,
			RepoListing: github.RepositoryListing{
				Users: []string{"test-user"},
			},
			MergeTypes: []scm.MergeType{scm.MergeTypeMerge},
		})
		require.NoError(t, err)

		repos, err := gh.GetRepositories(context.Background())
		assert.NoError(t, err)
		if assert.Len(t, repos, 1) {
			assert.Equal(t, "main", repos[0].DefaultBranch())
			assert.Equal(t, "lindell/test2", repos[0].FullName())
		}
	}

	// Topics
	{
		gh, err := github.New(github.Config{
			TransportMiddleware: transport.Wrapper,
			RepoListing: github.RepositoryListing{
				Organizations: []string{"test-org"},
				Topics:        []string{"frontend", "backend"},
			},
			MergeTypes: []scm.MergeType{scm.MergeTypeMerge},
		})
		require.NoError(t, err)

		repos, err := gh.GetRepositories(context.Background())
		assert.NoError(t, err)
		if assert.Len(t, repos, 1) {
			assert.Equal(t, "master", repos[0].DefaultBranch())
			assert.Equal(t, "test-org/test1", repos[0].FullName())
		}
	}

	// Multiple
	{
		gh, err := github.New(github.Config{
			TransportMiddleware: transport.Wrapper,
			RepoListing: github.RepositoryListing{
				Organizations: []string{"test-org"},
				Users:         []string{"test-user"},
				Repositories: []github.RepositoryReference{
					{
						OwnerName: "test-org",
						Name:      "test1",
					},
				},
			},
			MergeTypes: []scm.MergeType{scm.MergeTypeMerge},
		})
		require.NoError(t, err)

		repos, err := gh.GetRepositories(context.Background())
		assert.NoError(t, err)
		if assert.Len(t, repos, 2) {
			assert.Equal(t, "master", repos[0].DefaultBranch())
			assert.Equal(t, "test-org/test1", repos[0].FullName())
			assert.Equal(t, "main", repos[1].DefaultBranch())
			assert.Equal(t, "lindell/test2", repos[1].FullName())
		}
	}

	// Forks
	{
		gh, err := github.New(github.Config{
			TransportMiddleware: transport.Wrapper,
			RepoListing: github.RepositoryListing{
				Organizations: []string{"test-org"},
				Users:         []string{"test-user"},
				Repositories: []github.RepositoryReference{
					{
						OwnerName: "test-org",
						Name:      "test1",
					},
				},
				SkipForks: true,
			},
			MergeTypes: []scm.MergeType{scm.MergeTypeMerge},
		})
		require.NoError(t, err)

		repos, err := gh.GetRepositories(context.Background())
		assert.NoError(t, err)
		if assert.Len(t, repos, 1) {
			assert.Equal(t, "test-org/test1", repos[0].FullName())
		}
	}
}
