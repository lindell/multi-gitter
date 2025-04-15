package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/lindell/multi-gitter/internal/git"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
)

// Github should implement the ChangePusher interface
var _ scm.ChangePusher = &Github{}

func (g *Github) Push(
	ctx context.Context,
	r scm.Repository,
	changes []git.Changes,
	featureBranch string,
	branchExist bool,
	forcePush bool,
) error {
	repo := r.(repository)

	// There is no way to force push with the API, so we need to delete the branch
	// and create it again.
	if forcePush {
		err := g.deleteRef(ctx, repo, featureBranch)
		if err != nil {
			return err
		}

		branchExist = false
	}

	if !branchExist {
		err := g.CreateBranch(ctx, repo, featureBranch, changes[0].OldHash)
		if err != nil {
			return err
		}
	}

	var err error
	newHash := ""
	for _, change := range changes {
		// If multiple changes are made, the old hash should be the new hash
		// from the previous commit. The commit locally won't be exactly the same
		// as the commit made through the API
		if newHash != "" {
			change.OldHash = newHash
		}

		newHash, err = g.CommitThroughAPI(ctx, repo, featureBranch, change)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Github) CommitThroughAPI(ctx context.Context,
	repo repository,
	branch string,
	changes git.Changes,
) (string, error) {
	query := `
		mutation ($input: CreateCommitOnBranchInput!) {
			createCommitOnBranch(input: $input) {
				commit {
					oid
				}
			}
		}`

	var v createCommitOnBranchInput

	v.Input.Branch.RepositoryNameWithOwner = repo.ownerName + "/" + repo.name

	v.Input.Branch.BranchName = branch
	v.Input.ExpectedHeadOid = changes.OldHash
	v.Input.Message.Headline = changes.Message

	for path, contents := range changes.Additions {
		v.Input.FileChanges.Additions = append(v.Input.FileChanges.Additions, commitAddition{
			Path:     path,
			Contents: base64.StdEncoding.EncodeToString(contents),
		})
	}

	for _, path := range changes.Deletions {
		v.Input.FileChanges.Deletions = append(v.Input.FileChanges.Deletions, commitDeletion{
			Path: path,
		})
	}

	var result createCommitOnBranchOutput
	err := g.makeGraphQLRequest(ctx, query, v, &result)
	if err != nil {
		return "", errors.WithMessage(err, "could not commit changes though API")
	}
	oid := result.CreateCommitOnBranch.Commit.Oid
	if oid == "" {
		return "", errors.New("could not get commit oid")
	}

	return oid, nil
}

func (g *Github) CreateBranch(ctx context.Context, repo repository, branchName string, oid string) error {
	query := `mutation($input: CreateRefInput!){
		createRef(input: $input) {
			ref {
				name
			}
		}
	}`

	if !strings.HasPrefix(repo.name, "refs/heads/") {
		branchName = "refs/heads/" + branchName
	}

	var cri CreateRefInput
	cri.Input.Name = branchName

	cri.Input.Oid = oid
	cri.Input.RepositoryID = repo.id

	err := g.makeGraphQLRequest(ctx, query, cri, nil)
	if err != nil {
		return errors.WithMessage(err, "could not create branch")
	}

	return nil
}

func (g *Github) deleteRef(ctx context.Context, repo repository, branchName string) error {
	branchRef, err := g.getBranchID(ctx, repo, branchName)
	if err != nil {
		return err
	}

	query := `mutation($input: DeleteRefInput!){
		deleteRef(input: $input){
			clientMutationId
		}
	}`

	var deleteRefInput DeleteRefInput
	deleteRefInput.Input.RefID = branchRef

	err = g.makeGraphQLRequest(ctx, query, deleteRefInput, nil)
	if err != nil {
		return errors.WithMessage(err, "could not delete branch")
	}

	return nil
}

func (g *Github) getBranchID(ctx context.Context, repo repository, branchName string) (string, error) {
	query := `query($owner: String!, $name: String!, $qualifiedName: String!) {
		repository(owner: $owner, name: $name) {
			ref(qualifiedName: $qualifiedName) {
				id
			}
		}
	}`

	var result getRefOutput
	err := g.makeGraphQLRequest(ctx, query, &RepositoryInput{
		Name:          repo.name,
		Owner:         repo.ownerName,
		QualifiedName: fmt.Sprintf("refs/heads/%s", branchName),
	}, &result)
	if err != nil {
		return "", errors.WithMessage(err, "could not get branch ID")
	}

	return result.Repository.Ref.ID, nil
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
			Additions []commitAddition `json:"additions,omitempty"`
			Deletions []commitDeletion `json:"deletions,omitempty"`
		} `json:"fileChanges"`
	} `json:"input"`
}

type createCommitOnBranchOutput struct {
	CreateCommitOnBranch struct {
		Commit struct {
			Oid string `json:"oid"`
		} `json:"commit"`
	} `json:"createCommitOnBranch"`
}

type commitAddition struct {
	Path     string `json:"path,omitempty"`
	Contents string `json:"contents,omitempty"`
}

type commitDeletion struct {
	Path string `json:"path,omitempty"`
}

type CreateRefInput struct {
	Input struct {
		Name         string `json:"name"`
		Oid          string `json:"oid"`
		RepositoryID string `json:"repositoryId"`
	} `json:"input"`
}

type RepositoryInput struct {
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	QualifiedName string `json:"qualifiedName"`
}

type DeleteRefInput struct {
	Input struct {
		RefID string `json:"refId"`
	} `json:"input"`
}

type getRepositoryOutput struct {
	Repository struct {
		ID string `json:"id"`
	} `json:"repository"`
}

type getRefOutput struct {
	Repository struct {
		Ref struct {
			ID string `json:"id"`
		} `json:"ref"`
	} `json:"repository"`
}
