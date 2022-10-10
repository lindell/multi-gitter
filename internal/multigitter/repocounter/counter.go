package repocounter

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	mgerrors "github.com/lindell/multi-gitter/internal/multigitter/errors"
	"github.com/lindell/multi-gitter/internal/multigitter/terminal"
	"github.com/lindell/multi-gitter/internal/scm"
)

const QuestionCtrlC = -1

const descriptionLen = 11
const maxRepoNameLen = 50

// Counter keeps track of succeeded and failed repositories
type Counter struct {
	// Used for quick lookups
	repositoriesMap map[string]*repoStatus
	// Used for ordered listing
	repositoriesList []*repoStatus
	repoMaxLength    int

	columnsLengths []int

	done int // Either successful or faulty runs

	lock sync.RWMutex

	question     question
	questionLock sync.Mutex // Used

	screen tcell.Screen
}

type question struct {
	text    string
	index   int
	options []QuestionOption
}

type QuestionOption struct {
	Text     string
	Shortcut rune
}

type repoStatus struct {
	scm.Repository

	pullRequest scm.PullRequest
	err         error
	action      Action
}

// NewCounter create a new repo counter
func NewCounter(repos []scm.Repository) *Counter {
	counter := &Counter{
		repositoriesMap: map[string]*repoStatus{},
	}

	maxLength := 0
	for _, repo := range repos {
		status := &repoStatus{
			Repository: repo,
		}
		counter.repositoriesMap[repo.FullName()] = status
		counter.repositoriesList = append(counter.repositoriesList, status)

		nameLength := len(repo.FullName())
		if nameLength > maxLength {
			maxLength = nameLength
		}
	}

	if maxLength > maxRepoNameLen {
		maxLength = maxRepoNameLen
	}

	counter.repoMaxLength = maxLength
	counter.columnsLengths = make([]int, 3)
	counter.columnsLengths[0] = maxLength
	counter.columnsLengths[1] = descriptionLen

	return counter
}

func (r *Counter) columnStart(columnI int) int {
	start := 1
	for i := 0; i < columnI; i++ {
		start += r.columnsLengths[i] + 1
	}
	return start
}

func (r *Counter) OpenTTY() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}

	r.screen = s

	return nil
}

func (r *Counter) CloseTTY() error {
	r.screen.Fini()
	return nil
}

// AddError add a failing repository together with the error that caused it
func (r *Counter) SetError(err error, repo scm.Repository) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.repositoriesMap[repo.FullName()].action = ActionError
	r.repositoriesMap[repo.FullName()].err = err
	r.done++
	r.ttyRender()
}

func (r *Counter) SetRepoAction(repo scm.Repository, action Action) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.repositoriesMap[repo.FullName()].action = action
	if action == ActionSuccess || action == ActionError {
		r.done++
	}

	r.ttyRender()
}

// AddSuccessPullRequest adds a pullrequest that succeeded
func (r *Counter) AddSuccessPullRequest(repo scm.Repository, pr scm.PullRequest) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.repositoriesMap[repo.FullName()].pullRequest = pr
}

func (r *Counter) QuestionLock() {
	r.questionLock.Lock()
}

func (r *Counter) QuestionUnlock() {
	r.questionLock.Unlock()
}

func (r *Counter) AskQuestion(text string, options []QuestionOption) int {
	defer func() {
		r.question = question{}
		r.ttyRender()
	}()

	r.question = question{
		text:    text,
		options: options,
	}

	for {
		r.ttyRender()
		event := r.screen.PollEvent()
		if event, ok := event.(*tcell.EventKey); ok {
			if event.Key() == tcell.KeyCtrlC {
				return QuestionCtrlC
			}

			runeKey := event.Rune()
			if runeKey != 0 {
				for i, opt := range r.question.options {
					if opt.Shortcut == runeKey {
						return i
					}
				}
			}

			switch event.Key() {
			case tcell.KeyRight:
				r.question.index = (r.question.index + 1) % len(r.question.options)
			case tcell.KeyLeft:
				r.question.index--
				if r.question.index < 0 {
					r.question.index += len(r.question.options)
				}
			case tcell.KeyEnter:
				return r.question.index
			}
		}
	}
}

func (r *Counter) SuspendTTY() {
	_ = r.screen.Suspend()
}

func (r *Counter) ResumeTTY() {
	_ = r.screen.Resume()
}

func (r *Counter) ttyRender() {
	if r.screen == nil {
		return
	}

	r.screen.Clear()

	screenWidth, screenHeight := r.screen.Size()

	// Header
	headerStyle := tcell.StyleDefault.Background(tcell.ColorLightGrey).Foreground(tcell.ColorBlack)
	for i := 0; i < screenWidth; i++ {
		r.screen.SetContent(i, 0, ' ', nil, headerStyle)
	}
	emitStr(r.screen, r.columnStart(0), 0, headerStyle, "NAME")
	emitStr(r.screen, r.columnStart(1), 0, headerStyle, center("STATE", descriptionLen))
	emitStr(r.screen, r.columnStart(2), 0, headerStyle, "PROGRESS")

	y := 1

	for _, repo := range r.repositoriesList {
		description := repo.action.Description()
		if repo.action == ActionError {
			if repo.err == mgerrors.ErrNoChange {
				description = "no change"
			} else if repo.err == mgerrors.ErrRejected {
				description = "rejected"
			} else if err, ok := repo.err.(*exec.ExitError); ok {
				description = fmt.Sprintf("exit %d", err.ExitCode())
			}
		}

		emitStr(r.screen, r.columnStart(0), y, tcell.StyleDefault, shortenRepoName(repo, maxRepoNameLen))
		emitStr(r.screen, r.columnStart(1), y, repo.action.Color(), center(description, descriptionLen))
		progressBar(
			r.screen,
			// screenWidth-(r.repoMaxLength+descriptionLen+2),
			20,
			statusToPercentage(repo.action),
			r.columnStart(2),
			y,
		)

		y++
	}

	progressBarWithCounter(
		r.screen,
		screenWidth,
		len(r.repositoriesList),
		r.done,
		0,
		screenHeight-1,
	)

	// Not thread safe yet TODO fix
	emitStr(r.screen, 1, screenHeight-4, tcell.StyleDefault, r.question.text)
	buttonX := 1
	for i, question := range r.question.options {
		buttonX += button(r.screen, buttonX, screenHeight-3, question.Text, i == r.question.index) + 1
	}

	r.screen.Show()
}

// Info returns a formated string about all repositories
func (r *Counter) Info() string {
	defer r.lock.RUnlock()
	r.lock.RLock()

	var exitInfo string

	// Group all error messages together
	errMap := map[string][]*repoStatus{}
	for _, repo := range r.repositoriesList {
		if repo.err != nil {
			errMap[repo.err.Error()] = append(errMap[repo.err.Error()], repo)
		}
	}

	for errMsg := range errMap {
		exitInfo += fmt.Sprintf("%s:\n", strings.ToUpper(errMsg[0:1])+errMsg[1:])
		for _, repo := range errMap[errMsg] {
			exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
		}
	}

	successRepositories := []*repoStatus{}
	for _, repo := range r.repositoriesList {
		if repo.action == ActionSuccess {
			successRepositories = append(successRepositories, repo)
		}
	}

	if len(successRepositories) > 0 {
		exitInfo += "Repositories with a successful run:\n"
		for _, repoStatus := range successRepositories {
			if repoStatus.pullRequest != nil {
				if urler, ok := repoStatus.pullRequest.(urler); ok {
					exitInfo += fmt.Sprintf("  %s\n", terminal.Link(repoStatus.pullRequest.String(), urler.URL()))
				} else {
					exitInfo += fmt.Sprintf("  %s\n", repoStatus.pullRequest.String())
				}
			} else {
				exitInfo += fmt.Sprintf("  %s\n", repoStatus.FullName())
			}
		}
	}

	return exitInfo
}

// TTYSupported checks if TTY is supported
func TTYSupported() bool {
	_, err := tcell.NewScreen()
	return err == nil
}

type urler interface {
	URL() string
}
