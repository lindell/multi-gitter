package gerrit

import (
	gogerrit "github.com/andygrunwald/go-gerrit"
	"golang.org/x/net/context"
)

// The gogerrit.Client struct encapsulates sub-structs for each API endpoint.
// Defining our own interface allows us to flatten the client structure
// by exposing only the necessary methods, which simplifies mocking in tests.
type goGerritClient struct {
	client *gogerrit.Client
}

type GoGerritClient interface {
	ListProjects(ctx context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error)
	QueryChanges(ctx context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error)
	AbandonChange(ctx context.Context, changeID string, input *gogerrit.AbandonInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error)
	SubmitChange(ctx context.Context, changeID string, input *gogerrit.SubmitInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error)
}

func (ggc goGerritClient) ListProjects(ctx context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
	return ggc.client.Projects.ListProjects(ctx, opt)
}

func (ggc goGerritClient) QueryChanges(ctx context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
	return ggc.client.Changes.QueryChanges(ctx, opt)
}

func (ggc goGerritClient) AbandonChange(ctx context.Context, changeID string, input *gogerrit.AbandonInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error) {
	return ggc.client.Changes.AbandonChange(ctx, changeID, input)
}

func (ggc goGerritClient) SubmitChange(ctx context.Context, changeID string, input *gogerrit.SubmitInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error) {
	return ggc.client.Changes.SubmitChange(ctx, changeID, input)
}
