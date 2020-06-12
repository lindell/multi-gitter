package multigitter

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"

	"github.com/lindell/multi-gitter/internal/domain"
	"github.com/lindell/multi-gitter/internal/git"
)

// RepoGetter fetches repositories
type RepoGetter interface {
	GetRepositories() ([]domain.Repository, error)
}

// PullRequestCreator creates pull requests
type PullRequestCreator interface {
	CreatePullRequest(repo domain.Repository, newPR domain.NewPullRequest) error
}

// Runner conains fields to be able to do the run
type Runner struct {
	ScriptPath    string // Must be absolute path
	FeatureBranch string

	RepoGetter         RepoGetter
	PullRequestCreator PullRequestCreator

	CommitMessage    string
	PullRequestTitle string
	PullRequestBody  string
	Reviewers        []string
	MaxReviewers     int // If set to zero, all reviewers will be used
}

// Run runs a script for multiple repositories and creates PRs with the changes made
func (r Runner) Run() error {
	repos, err := r.RepoGetter.GetRepositories()
	if err != nil {
		return err
	}

	log.Printf("Running on %d repositories\n", len(repos))

	exitCodeRepos := map[int][]domain.Repository{}
	noChangeRepos := []domain.Repository{}
	successRepos := []domain.Repository{}

	defer func() {
		var exitInfo string
		for exitCode := range exitCodeRepos {
			exitInfo += fmt.Sprintf("Repositories with exit code %d:\n", exitCode)
			for _, repo := range exitCodeRepos[exitCode] {
				exitInfo += fmt.Sprintf("  %s\n", repo.GetURL())
			}
		}

		if len(noChangeRepos) > 0 {
			exitInfo += "Repositories where nothing was changed:\n"
			for _, repo := range noChangeRepos {
				exitInfo += fmt.Sprintf("  %s\n", repo.GetURL())
			}
		}

		if len(successRepos) > 0 {
			exitInfo += "Repositories with a successful run:\n"
			for _, repo := range successRepos {
				exitInfo += fmt.Sprintf("  %s\n", repo.GetURL())
			}
		}

		if exitInfo != "" {
			log.Print(exitInfo)
		}
	}()

	for _, repo := range repos {
		log.Printf("Cloning and running script on: %s\n", repo.GetURL())
		err := r.runSingleRepo(repo.GetURL())
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("Got exit code %d when running %s\n", exitErr.ExitCode(), repo.GetURL())
			exitCodeRepos[exitErr.ExitCode()] = append(exitCodeRepos[exitErr.ExitCode()], repo)
			continue
		} else if err == domain.NoChangeError {
			log.Printf("No change done on the repo by the script when running: %s\n", repo.GetURL())
			noChangeRepos = append(noChangeRepos, repo)
			continue
		} else if err != nil {
			return err
		}

		err = r.PullRequestCreator.CreatePullRequest(repo, domain.NewPullRequest{
			Title:     r.PullRequestTitle,
			Body:      r.PullRequestBody,
			Head:      r.FeatureBranch,
			Base:      repo.GetBranch(),
			Reviewers: getReviewers(r.Reviewers, r.MaxReviewers),
		})
		if err != nil {
			return err
		}

		successRepos = append(successRepos, repo)
		return nil
	}

	return nil
}

func getReviewers(reviewers []string, maxReviewers int) []string {
	if maxReviewers == 0 || len(reviewers) <= maxReviewers {
		return reviewers
	}

	rand.Shuffle(len(reviewers), func(i, j int) { reviewers[i], reviewers[j] = reviewers[j], reviewers[i] })

	return reviewers[0:maxReviewers]
}

func (r Runner) runSingleRepo(url string) error {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-changer-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	sourceController := git.Git{
		Directory: tmpDir,
		Repo:      url,
		NewBranch: r.FeatureBranch,
	}

	err = sourceController.Clone()
	if err != nil {
		return err
	}

	// Run the command that might or might not change the content of the repo
	// If the command return a non zero exit code, abort.
	cmd := exec.Command(r.ScriptPath)
	cmd.Dir = tmpDir

	reader, writer := io.Pipe()
	cmd.Stdout = writer
	cmd.Stderr = writer

	// Print each line that is outputted by the script
	go func() {
		buf := bufio.NewReader(reader)
		for {
			line, err := buf.ReadString('\n')
			if line != "" {
				log.Printf("Script output: %s", line)
			}
			if err != nil {
				return
			}
		}
	}()

	err = cmd.Run()
	if err != nil {
		return err
	}

	err = sourceController.Commit(r.CommitMessage)
	if err != nil {
		return err
	}

	return nil
}
