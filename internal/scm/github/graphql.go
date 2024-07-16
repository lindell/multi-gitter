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

func (g *Github) getRepositoryID(ctx context.Context, owner string, name string) (string, error) {
	query := `query($owner: String!, $name: String!) {
		repository(owner: $owner, name: $name) {
			id
		} 
	}`

	var v getRepositoryInput
	var result getRepositoryOutput
	v.Name = name
	v.Owner = owner

	err := g.makeGraphQLRequest(ctx, query, v, &result)

	if err != nil {
		return "", err
	}

	return result.Repository.ID, nil
}

func (g *Github) createRef(ctx context.Context, branchName string, oid string, iD string) error {
	query := `mutation($input: CreateRefInput!){
		createRef(input: $input) {
			ref {
				name
			}
		}
	}`

	var v createRefInput

	v.Input.Name = "refs/heads/" + branchName
	v.Input.Oid = oid
	v.Input.RepositoryID = iD

	fmt.Printf("%v", v)

	var result map[string]interface{}

	err := g.makeGraphQLRequest(ctx, query, v, &result)

	if err != nil {
		return err
	}

	return nil
}

func (g *Github) CommitThroughAPI(ctx context.Context, input CreateCommitOnBranchInput) error {
	query := `
		mutation ($input: CreateCommitOnBranchInput!) {
			createCommitOnBranch(input: $input) {
				commit {
				url
				}
			}
		}`

	var v createCommitOnBranchInput

	repoID, err := g.getRepositoryID(ctx, input.Owner, input.RepositoryName)

	if err != nil {
		return err
	}

	err = g.createRef(ctx, input.BranchName, input.ExpectedHeadOid, repoID)

	if err != nil {
		return err
	}

	v.Input.Branch.RepositoryNameWithOwner = input.Owner + "/" + input.RepositoryName

	v.Input.Branch.BranchName = input.BranchName
	v.Input.ExpectedHeadOid = input.ExpectedHeadOid
	v.Input.Message.Headline = input.Headline

	for path, contents := range input.Additions {
		v.Input.FileChanges.Additions = append(v.Input.FileChanges.Additions, struct {
			Path     string "json:\"path,omitempty\""
			Contents string "json:\"contents,omitempty\""
		}{Path: path, Contents: contents})
	}

	for _, path := range input.Deletions {
		v.Input.FileChanges.Deletions = append(v.Input.FileChanges.Deletions, struct {
			Path string "json:\"path,omitempty\""
		}{Path: path})
	}

	var result map[string]interface{}

	err = g.makeGraphQLRequest(ctx, query, v, &result)

	if err != nil {
		return err
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

type createCommitOnBranchInput struct {
	Input struct {
		ExpectedHeadOid string `json:"expectedHeadOid"`
		Branch          struct {
			RepositoryNameWithOwner string `json:"repositoryNameWithOwner"`
			BranchName              string `json:"branchName"`
		} `json:"branch"`
		Message struct {
			Headline string `json:"headline"`
		} `json:"message"`
		FileChanges struct {
			Additions []struct {
				Path     string `json:"path,omitempty"`
				Contents string `json:"contents,omitempty"`
			} `json:"additions"`
			Deletions []struct {
				Path string `json:"path,omitempty"`
			} `json:"deletions"`
		} `json:"fileChanges"`
	} `json:"input"`
}

type createRefInput struct {
	Input struct {
		Name         string `json:"name"`
		Oid          string `json:"oid"`
		RepositoryID string `json:"repositoryId"`
	} `json:"input"`
}

type getRepositoryInput struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

type getRepositoryOutput struct {
	Repository struct {
		ID string `json:"id"`
	} `json:"repository"`
}

type CreateCommitOnBranchInput struct {
	RepositoryName  string
	Owner           string
	BranchName      string
	Headline        string
	Additions       map[string]string
	Deletions       []string
	ExpectedHeadOid string
}
