package repocounter_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/lindell/multi-gitter/tests/vcmock"
)

func fakeRepo(i int) vcmock.Repository {
	return vcmock.Repository{
		OwnerName: fmt.Sprintf("owner-%d", i),
		RepoName:  fmt.Sprintf("repo-%d", i),
	}
}

func fakePR(i int) vcmock.PullRequest {
	return vcmock.PullRequest{
		PRStatus:   scm.PullRequestStatusPending,
		PRNumber:   42,
		Merged:     false,
		Repository: fakeRepo(i),
	}
}

func TestCounter_Info(t *testing.T) {
	tests := []struct {
		name     string
		changeFn func(r *repocounter.Counter)
		want     string
	}{
		{
			name:     "empty",
			want:     "",
			changeFn: func(r *repocounter.Counter) {},
		},
		{
			name: "one success",
			changeFn: func(r *repocounter.Counter) {
				r.AddSuccessRepositories(fakeRepo(1))
			},
			want: `
Repositories with a successful run:
  owner-1/repo-1
`,
		},
		{
			name: "one success with pr",
			changeFn: func(r *repocounter.Counter) {
				r.AddSuccessPullRequest(fakeRepo(1), fakePR(1))
			},
			want: `
Repositories with a successful run:
  owner-1/repo-1 #42
`,
		},
		{
			name: "one error no pr",
			changeFn: func(r *repocounter.Counter) {
				r.AddError(errors.New("test error"), fakeRepo(1), nil)
			},
			want: `
Test error:
  owner-1/repo-1
`,
		},
		{
			name: "one error with pr",
			changeFn: func(r *repocounter.Counter) {
				r.AddError(errors.New("test error"), fakeRepo(1), fakePR(1))
			},
			want: `
Test error:
  owner-1/repo-1 #42
`,
		},
		{
			name: "one error with pr",
			changeFn: func(r *repocounter.Counter) {
				r.AddError(errors.New("test error"), fakeRepo(1), fakePR(1))
			},
			want: `
Test error:
  owner-1/repo-1 #42
`,
		},
		{
			name: "multiple",
			changeFn: func(r *repocounter.Counter) {
				r.AddError(errors.New("test error"), fakeRepo(1), fakePR(1))
				r.AddSuccessPullRequest(fakeRepo(2), fakePR(2))
				r.AddSuccessPullRequest(fakeRepo(3), fakePR(3))
				r.AddError(errors.New("test error 2"), fakeRepo(5), fakePR(5))
				r.AddError(errors.New("test error"), fakeRepo(4), fakePR(4))
			},
			want: `
Test error:
  owner-1/repo-1 #42
  owner-4/repo-4 #42
Test error 2:
  owner-5/repo-5 #42
Repositories with a successful run:
  owner-2/repo-2 #42
  owner-3/repo-3 #42
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repocounter.NewCounter()
			tt.changeFn(r)

			if got := r.Info(); got != strings.TrimLeft(tt.want, "\n") {
				t.Errorf("Counter.Info() = %v, want %v", got, tt.want)
			}
		})
	}
}
