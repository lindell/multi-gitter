package gerrit

import (
	"context"
	"crypto/sha1" // #nosec
	"encoding/hex"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"time"

	gogerrit "github.com/andygrunwald/go-gerrit"
	internalHTTP "github.com/lindell/multi-gitter/internal/http"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
)

const FooterBranch = "MultiGitter-Branch"
const FooterChangeID = "Change-Id"
const QueryChangesLimit = 100
const QueryProjectsLimit = 100

type Gerrit struct {
	client     GoGerritClient
	baseURL    string
	username   string
	token      string
	repoSearch string
}

func New(username, token, baseURL, repoSearch string) (*Gerrit, error) {
	ctx := context.Background() // cancellation won't happen in our case (only used by go-gerrit if you inject username and token directly within baseURL)
	client, err := gogerrit.NewClient(ctx, baseURL, &http.Client{
		Transport: internalHTTP.LoggingRoundTripper{},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gerrit client")
	}

	client.Authentication.SetBasicAuth(username, token)

	return &Gerrit{
		client: goGerritClient{
			client: client,
		},
		baseURL:    baseURL,
		username:   username,
		token:      token,
		repoSearch: repoSearch,
	}, nil
}

func (g Gerrit) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
	repos := make([]scm.Repository, 0)
	skip := 0
	for {
		opt := &gogerrit.ProjectOptions{
			Description: true,
			Regex:       g.repoSearch,
			Type:        "CODE",
			Skip:        strconv.Itoa(skip),
			ProjectBaseOptions: gogerrit.ProjectBaseOptions{
				Limit: QueryProjectsLimit,
			},
		}
		projects, _, err := g.client.ListProjects(ctx, opt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list projects")
		}

		moreProjects := false
		for name, project := range *projects {
			if project.State == "ACTIVE" {
				repo, err := g.convertRepo(name)
				if err != nil {
					return nil, err
				}

				repos = append(repos, repo)
			}
			// Because projects is a map, Golang iteration do not guarantee order of projects.
			// So we cannot rely only on the last project to determine if there are more projects.
			if project.MoreProjects {
				moreProjects = true
			}
		}

		if !moreProjects {
			break
		}
		skip += QueryProjectsLimit
	}

	// Keep consistent order of repositories
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].(repository).name < repos[j].(repository).name
	})
	return repos, nil
}

func (g Gerrit) convertRepo(name string) (repository, error) {
	// Note: maybe we should support cloning via ssh
	u, err := url.Parse(g.baseURL)
	if err != nil {
		return repository{}, err
	}
	u.User = url.UserPassword(g.username, g.token)
	u.Path = "/a/" + name
	repoURL := u.String()

	return repository{
		url:           repoURL,
		name:          name,
		defaultBranch: "master", // Some projects might have a different default branch
	}, nil
}

func (g Gerrit) CreatePullRequest(ctx context.Context, repo scm.Repository, _ scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	// In Gerrit context, pushing a commit to refs/for/<base_branch> is enough to create automatically a change.
	// So here, we are just "fetching" the change related to current branch (Head of PR)
	// Not yet implemented: reviewers, team reviewers, assignees, draft, labels

	return g.getChange(ctx, repo, newPR.Head)
}

func (g Gerrit) UpdatePullRequest(ctx context.Context, repo scm.Repository, _ scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	// In Gerrit context, pushing a commit to refs/for/<base_branch> is enough to create automatically a change.
	// So here, we are just "fetching" the change related
	// Not yet implemented: reviewers, team reviewers, assignees, draft, labels

	return g.getChange(ctx, repo, updatedPR.Head)
}

func (g Gerrit) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
	repositories, err := g.GetRepositories(ctx)
	if err != nil {
		return nil, err
	}

	// Build a map of repository names to fast search if a change belongs to a repository
	projectNames := make(map[string]struct{}, len(repositories))
	for _, s := range repositories {
		projectNames[s.(repository).name] = struct{}{}
	}

	var prs []scm.PullRequest
	start := 0
	for {
		// Query all changes related to the branch name to avoid one query per repository
		changes, err := g.queryChanges(ctx, branchName, []string{}, start, QueryChangesLimit)
		if err != nil {
			return nil, err
		}

		for _, change := range changes {
			if _, ok := projectNames[change.Project]; ok {
				prs = append(prs, convertChange(change, g.baseURL))
			}
		}

		// The lastest change retuned by API indicates if there are more changes to fetch.
		if !changes[len(changes)-1].MoreChanges {
			break
		}
		start += QueryChangesLimit
	}

	// Keep consistent order of PRs
	sort.Slice(prs, func(i, j int) bool {
		return prs[i].(change).project < prs[j].(change).project
	})
	return prs, err
}

func (g Gerrit) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	changes, err := g.queryChanges(ctx, branchName, []string{
		"project:" + repo.FullName(),
		"is:open",
	}, 0, 5) // Limit to few changes, since we only care about the first one
	if err != nil {
		return nil, err
	}

	if len(changes) == 0 {
		return nil, nil
	} else if len(changes) > 1 {
		return nil, errors.New("More than one open change for branch " + branchName + " in project " + repo.FullName())
	}

	return convertChange(changes[0], g.baseURL), nil
}

func (g Gerrit) MergePullRequest(ctx context.Context, pr scm.PullRequest) error {
	change := pr.(change)

	_, _, err := g.client.SubmitChange(ctx, change.id, &gogerrit.SubmitInput{})

	return err
}

func (g Gerrit) ClosePullRequest(ctx context.Context, pr scm.PullRequest) error {
	change := pr.(change)

	_, _, err := g.client.AbandonChange(ctx, change.id, &gogerrit.AbandonInput{})
	if err != nil {
		return err
	}
	return nil
}

func (Gerrit) ForkRepository(_ context.Context, _ scm.Repository, _ string) (scm.Repository, error) {
	return nil, errors.New("Forking repositories is not supported in Gerrit")
}

func (g Gerrit) getChange(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	pr, err := g.GetOpenPullRequest(ctx, repo, branchName)
	if err != nil {
		return nil, err
	} else if pr == nil {
		return nil, errors.Errorf("Unable to find any open change related to branch %s in project %s", branchName, repo.FullName())
	}
	return pr, nil
}

func (g Gerrit) queryChanges(ctx context.Context, branchName string, filters []string, start int, limit int) ([]gogerrit.ChangeInfo, error) {
	defaultFilters := []string{
		"footer:" + FooterBranch + "=" + branchName,
	}
	query := strings.Join(append(defaultFilters, filters...), "+")

	opt := &gogerrit.QueryChangeOptions{
		QueryOptions: gogerrit.QueryOptions{
			Query: []string{query},
			Start: start,
			Limit: limit,
		},
		ChangeOptions: gogerrit.ChangeOptions{
			AdditionalFields: []string{
				"SUBMITTABLE",
			},
		},
	}
	changes, _, err := g.client.QueryChanges(ctx, opt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query changes: '%s'", filters)
	}
	return *changes, nil
}

func (g Gerrit) EnhanceCommit(ctx context.Context, repo scm.Repository, branchName string, commitMessage string) (string, error) {
	pr, err := g.GetOpenPullRequest(ctx, repo, branchName)
	if err != nil {
		return commitMessage, err
	}

	changeID := ""
	if pr != nil {
		changeID = pr.(change).changeID
	} else {
		changeID = generateChangeID(commitMessage)
	}
	message := commitMessage
	message += "\n\n" + FooterBranch + ": " + branchName
	message += "\n" + FooterChangeID + ": " + changeID
	return message, nil
}

func (g Gerrit) FeatureBranchExist(ctx context.Context, repo scm.Repository, branchName string) (bool, error) {
	pr, err := g.GetOpenPullRequest(ctx, repo, branchName)
	return pr != nil, err
}

func (g Gerrit) RemoteReference(baseBranch string, featureBranch string, skipPullRequest bool, pushOnly bool) string {
	if !skipPullRequest && !pushOnly {
		return "refs/for/" + baseBranch
	}
	return "refs/heads/" + featureBranch
}

func generateChangeID(commitMessage string) string {
	h := sha1.New() // #nosec
	hostname, _ := os.Hostname()
	whoami, _ := user.Current()
	h.Write([]byte(hostname))
	h.Write([]byte(whoami.Username))
	h.Write([]byte(time.Now().String()))
	h.Write([]byte(commitMessage))

	return "I" + hex.EncodeToString(h.Sum(nil))
}
