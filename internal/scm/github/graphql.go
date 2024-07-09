package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func (g *Github) makeGraphQLRequest(ctx context.Context, query string, data interface{}, res interface{}) error {
	type reqData struct {
		Query string      `json:"query"`
		Data  interface{} `json:"variables"`
	}
	rawReqData, err := json.Marshal(reqData{
		Query: query,
		Data:  data,
	})

	if err != nil {
		return errors.WithMessage(err, "could not marshal graphql request")
	}

	graphQLURL := "https://api.github.com/graphql"
	if g.baseURL != "" {
		graphQLURL, err = graphQLEndpoint(g.baseURL)
		if err != nil {
			return errors.WithMessage(err, "could not get graphql endpoint")
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", graphQLURL, bytes.NewBuffer(rawReqData))
	if err != nil {
		return errors.WithMessage(err, "could not create graphql request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.token))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resultData := struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
		Message string `json:"message"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&resultData); err != nil {
		return errors.WithMessage(err, "could not read graphql response body")
	}

	if len(resultData.Errors) > 0 {
		errorsMsgs := make([]string, len(resultData.Errors))
		for i := range resultData.Errors {
			errorsMsgs[i] = resultData.Errors[i].Message
		}
		return errors.WithMessage(
			errors.New(strings.Join(errorsMsgs, "\n")),
			"encountered error during GraphQL query",
		)
	}

	if resp.StatusCode >= 400 {
		return errors.Errorf("could not make GitHub GraphQL request: %s", resultData.Message)
	}

	if err := json.Unmarshal(resultData.Data, res); err != nil {
		return errors.WithMessage(err, "could not unmarshal graphQL result")
	}

	return nil
}

// graphQLEndpoint takes a url to a github enterprise instance (or the v3 api) and returns the url to the graphql endpoint
func graphQLEndpoint(u string) (string, error) {
	baseEndpoint, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(baseEndpoint.Path, "/") {
		baseEndpoint.Path += "/"
	}

	if strings.HasPrefix(baseEndpoint.Host, "api.") ||
		strings.Contains(baseEndpoint.Host, ".api.") {
		baseEndpoint.Path += "graphql"
	} else {
		baseEndpoint.Path = stripSuffixIfExist(baseEndpoint.Path, "v3/")
		baseEndpoint.Path = stripSuffixIfExist(baseEndpoint.Path, "api/")
		baseEndpoint.Path += "api/graphql"
	}

	return baseEndpoint.String(), nil
}

func (g *Github) graphQLCommitOnBranch(ctx context.Context) error {
	query := `
		mutation ($input: CreateCommitOnBranchInput!) {
			createCommitOnBranch(input: $input) {
				commit {
				url
				}
			}
		}
		`
	vars := `
	{
	"input": {
		"branch": {
			"repositoryNameWithOwner": "chrisstatham/multi-gitter",
			"branchName": "APPSEC-1108"
		},
		"message": {
			"headline": "Hello from GraphQL!"
		},
		"fileChanges": {
			"additions": [
				{
					"path": "myfile.txt",
					"contents": "SGVsbG8gZnJvbSBKQVZBIGFuZCBHcmFwaFFM"      
				}
			]
		},
		"expectedHeadOid": "cb65177542c9239b9773fb4f6e117c46108c32d2" 
		}
	}`

	result := map[string]graphqlRepo{}

	err := g.makeGraphQLRequest(ctx, query, vars, &result)
	if err != nil {
		fmt.Printf("FAil")
	}

	return nil
}

type graphqlPullRequestState string

const (
	graphqlPullRequestStateError   graphqlPullRequestState = "ERROR"
	graphqlPullRequestStateFailure graphqlPullRequestState = "FAILURE"
	graphqlPullRequestStatePending graphqlPullRequestState = "PENDING"
	graphqlPullRequestStateSuccess graphqlPullRequestState = "SUCCESS"
)

type graphqlRepo struct {
	PullRequests struct {
		Nodes []graphqlPR `json:"nodes"`
	} `json:"pullRequests"`
}

type graphqlPR struct {
	Number         int    `json:"number"`
	HeadRefName    string `json:"headRefName"`
	Closed         bool   `json:"closed"`
	URL            string `json:"url"`
	Merged         bool   `json:"merged"`
	BaseRepository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"baseRepository"`
	HeadRepository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"headRepository"`
	Commits struct {
		Nodes []struct {
			Commit struct {
				StatusCheckRollup struct {
					State *graphqlPullRequestState `json:"state"`
				} `json:"statusCheckRollup"`
			} `json:"commit"`
		} `json:"nodes"`
	} `json:"commits"`
}
