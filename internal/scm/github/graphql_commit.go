package github

import (
	"context"
	"fmt"
	"strings"
)

func (g *Github) CommitAndPushThoughGraphQL(ctx context.Context,
	headline string,
	featureBranch string,
	cloneURL string,
	oldHash string,
	additions map[string]string,
	deletions []string,
	forcePush bool,
	branchExist bool) error {
	array := strings.Split(cloneURL, "/")

	repositoryName := strings.Trim(array[len(array)-1], ".git")
	owner := array[len(array)-2]

	if forcePush {
		fmt.Printf("about to get branchid\n")

		branchID, err := g.getBranchID(ctx, owner, repositoryName, featureBranch)
		if err != nil {
			return err
		}

		fmt.Printf("about to deleteref\n")

		err = g.deleteRef(ctx, branchID)
		if err != nil {
			return err
		}

		branchExist = false
	}

	if !branchExist {
		fmt.Printf("about to create branch\n")

		err := g.CreateBranch(ctx, owner,
			repositoryName,
			featureBranch,
			oldHash)
		if err != nil {
			return err
		}
	}

	fmt.Printf("About to commit though API\n")

	err := g.CommitThroughAPI(ctx, owner, repositoryName, featureBranch, oldHash, headline, additions, deletions)

	if err != nil {
		return err
	}

	return nil
}

func (g *Github) getRepositoryID(ctx context.Context, owner string, name string) (string, error) {
	query := `query($owner: String!, $name: String!) {
		repository(owner: $owner, name: $name) {
			id
		} 
	}`

	var result getRepositoryOutput

	err := g.makeGraphQLRequest(ctx, query, &RepositoryInput{Name: name, Owner: owner}, &result)

	if err != nil {
		return "", err
	}

	return result.Repository.ID, nil
}

func (g *Github) getBranchID(ctx context.Context, owner string, repoName string, branchName string) (string, error) {
	query := `query($owner: String!, $name: String!) {
		repository(owner: $owner, name: $name) {
			refs(first: 100, refPrefix: "refs/heads/") {
				edges {
					node {
						id
						name
					}
				}
			}
		} 
	}`

	var result getRefsOutput

	err := g.makeGraphQLRequest(ctx, query, &RepositoryInput{Name: repoName, Owner: owner}, &result)

	if err != nil {
		return "", err
	}

	for _, edge := range result.Repository.Refs.Edges {
		if edge.Node.Name == branchName {
			return edge.Node.ID, nil
		}
	}

	return "", fmt.Errorf("unable to find branch named %s to delete", branchName)
}

func (g *Github) CreateBranch(ctx context.Context, owner string, repoName string, branchName string, oid string) error {
	query := `mutation($input: CreateRefInput!){
		createRef(input: $input) {
			ref {
				name
			}
		}
	}`

	var cri CreateRefInput

	repoID, err := g.getRepositoryID(ctx, owner, repoName)

	if err != nil {
		return err
	}

	if !strings.HasPrefix(repoName, "refs/heads/") {
		cri.Input.Name = "refs/heads/" + branchName
	} else {
		cri.Input.Name = branchName
	}

	cri.Input.Oid = oid
	cri.Input.RepositoryID = repoID

	var result interface{}

	err = g.makeGraphQLRequest(ctx, query, cri, &result)

	if err != nil {
		return err
	}

	return nil
}

func (g *Github) CommitThroughAPI(ctx context.Context, owner string, repoName string, branch string, oid string, headline string, additions map[string]string, deletions []string) error {
	query := `
		mutation ($input: CreateCommitOnBranchInput!) {
			createCommitOnBranch(input: $input) {
				commit {
				url
				}
			}
		}`

	var v createCommitOnBranchInput

	v.Input.Branch.RepositoryNameWithOwner = owner + "/" + repoName

	v.Input.Branch.BranchName = branch
	v.Input.ExpectedHeadOid = oid
	v.Input.Message.Headline = headline

	for path, contents := range additions {
		v.Input.FileChanges.Additions = append(v.Input.FileChanges.Additions, struct {
			Path     string "json:\"path,omitempty\""
			Contents string "json:\"contents,omitempty\""
		}{Path: path, Contents: contents})
	}

	for _, path := range deletions {
		v.Input.FileChanges.Deletions = append(v.Input.FileChanges.Deletions, struct {
			Path string "json:\"path,omitempty\""
		}{Path: path})
	}

	var result map[string]interface{}

	err := g.makeGraphQLRequest(ctx, query, v, &result)

	if err != nil {
		return err
	}

	return nil
}

func (g *Github) deleteRef(ctx context.Context, branchRef string) error {
	query := `mutation($input: DeleteRefInput!){
		deleteRef(input: $input){
    		clientMutationId
  		}
	}`

	var deleteRefInput DeleteRefInput
	deleteRefInput.Input.RefID = branchRef

	type ignoreReturn map[string]interface{}

	err := g.makeGraphQLRequest(ctx, query, deleteRefInput, &ignoreReturn{})

	if err != nil {
		return err
	}

	return nil
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

type CreateRefInput struct {
	Input struct {
		Name         string `json:"name"`
		Oid          string `json:"oid"`
		RepositoryID string `json:"repositoryId"`
	} `json:"input"`
}

type RepositoryInput struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
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

type getRefsOutput struct {
	Repository struct {
		Refs struct {
			Edges []struct {
				Node struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"refs"`
	} `json:"repository"`
}
