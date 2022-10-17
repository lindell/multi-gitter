package repocounter

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gdamore/tcell/v2"
	"github.com/lindell/go-ordered-set/orderedset"
	mgerrors "github.com/lindell/multi-gitter/internal/multigitter/errors"
	"github.com/lindell/multi-gitter/internal/multigitter/terminal"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/sirupsen/logrus"
)

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
	questionLock sync.Mutex

	toBeStarted *orderedset.OrderedSet[*repoStatus]
	inProgress  *orderedset.OrderedSet[*repoStatus]
	completed   *orderedset.OrderedSet[*repoStatus]

	screen            tcell.Screen
	quitCh            chan struct{}
	eventCallback     func(tcell.Event)
	eventCallbackLock sync.Mutex

	abortListeners     []func()
	abortListenersLock sync.Mutex

	lastRender      time.Time
	renderQueueLock sync.Mutex
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
		toBeStarted:     orderedset.New[*repoStatus](),
		inProgress:      orderedset.New[*repoStatus](),
		completed:       orderedset.New[*repoStatus](),
	}

	maxLength := 0
	for _, repo := range repos {
		status := &repoStatus{
			Repository: repo,
		}
		counter.repositoriesMap[repo.FullName()] = status
		counter.toBeStarted.Add(status)
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

	ch := make(chan tcell.Event)
	r.quitCh = make(chan struct{})
	go func() {
		s.ChannelEvents(ch, r.quitCh)
	}()
	go func() {
		for event := range ch {
			r.handleEvent(event)
		}
	}()

	r.screen = s

	return nil
}

func (r *Counter) CloseTTY() error {
	r.quitCh <- struct{}{}
	r.screen.Fini()
	return nil
}

func (r *Counter) OnAbort(fn func()) {
	r.abortListenersLock.Lock()
	defer r.abortListenersLock.Unlock()

	r.abortListeners = append(r.abortListeners, fn)
}

func (r *Counter) abort() {
	r.abortListenersLock.Lock()
	defer r.abortListenersLock.Unlock()

	for _, fn := range r.abortListeners {
		fn()
	}
}

// AddError add a failing repository together with the error that caused it
func (r *Counter) SetError(err error, repo scm.Repository) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.setRepoAction(repo, ActionError)
	r.repositoriesMap[repo.FullName()].err = err
	r.queueRender()
}

func (r *Counter) SetRepoAction(repo scm.Repository, action Action) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.setRepoAction(repo, action)

	r.queueRender()
}

func (r *Counter) waitForEvent() tcell.Event {
	var event tcell.Event
	var wg sync.WaitGroup
	wg.Add(1)

	r.eventCallbackLock.Lock()
	r.eventCallback = func(e tcell.Event) {
		event = e
		wg.Done()
	}
	r.eventCallbackLock.Unlock()

	wg.Wait()

	r.eventCallbackLock.Lock()
	r.eventCallback = nil
	r.eventCallbackLock.Unlock()

	return event
}

func (r *Counter) handleEvent(event tcell.Event) {
	logrus.Info(spew.Sdump(event))

	if event, ok := event.(*tcell.EventKey); ok {
		if event.Key() == tcell.KeyCtrlC {
			r.abort()
			return
		}
	}

	if _, ok := event.(*tcell.EventResize); ok {
		r.queueRender()
		return
	}

	r.eventCallbackLock.Lock()
	if r.eventCallback != nil {
		r.eventCallback(event)
	}
	r.eventCallbackLock.Unlock()
}

func (r *Counter) setRepoAction(repo scm.Repository, action Action) {
	status := r.repositoriesMap[repo.FullName()]
	status.action = action
	percentage := action.Percentage()
	if percentage == 1 {
		r.inProgress.Delete(status)
		r.completed.Add(status)
		r.done++
	} else if percentage > 0 {
		r.toBeStarted.Delete(status)
		r.inProgress.Add(status)
	}
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
	}()

	r.question = question{
		text:    text,
		options: options,
	}

	for {
		r.ttyRender()
		event := r.waitForEvent()
		if event, ok := event.(*tcell.EventKey); ok {
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

func (r *Counter) queueRender() {
	r.renderQueueLock.Lock()
	defer r.renderQueueLock.Unlock()
	// TODO: This is not working (as expected)

	sinceLastRender := time.Since(r.lastRender)
	if sinceLastRender > time.Millisecond*1 {
		r.ttyRender()
		return
	}
	time.Sleep(time.Millisecond*1 - sinceLastRender)
	r.ttyRender()
}

func (r *Counter) ttyRender() {
	if r.screen == nil {
		return
	}

	r.lastRender = time.Now()

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

	// Determine the size of each category, completed/in progress/to be started
	totalSize := screenHeight - 3
	inProgressSize := minInt(r.inProgress.Size(), totalSize) // In progress takes priority, but can't still be bigger than the screen
	completedSize := minInt((totalSize-inProgressSize)/2, r.completed.Size())
	toBeStartedSize := minInt(totalSize-inProgressSize-completedSize, r.toBeStarted.Size())
	completedSize = maxInt(completedSize, totalSize-inProgressSize-toBeStartedSize)

	y := 0
	completedY := 0
	iterateTimes[*repoStatus](completedSize, r.completed.IterReverse(), func(repo *repoStatus) {
		r.renderRepoStatus(y+completedSize-completedY, repo)
		completedY++
	})
	y += completedSize

	inProgressY := 0
	iterateTimes[*repoStatus](inProgressSize, r.inProgress.IterReverse(), func(repo *repoStatus) {
		r.renderRepoStatus(y+inProgressSize-inProgressY, repo)
		inProgressY++
	})
	y += inProgressSize

	iterateTimes[*repoStatus](toBeStartedSize, r.toBeStarted.Iter(), func(repo *repoStatus) {
		y++
		r.renderRepoStatus(y, repo)
	})

	progressBarWithCounter(
		r.screen,
		screenWidth,
		len(r.repositoriesList),
		r.done,
		0,
		screenHeight-1,
	)

	// Not thread safe yet TODO fix
	emitStr(r.screen, 1, screenHeight-3, headerStyle, r.question.text)
	buttonX := 1
	for i, question := range r.question.options {
		buttonX += button(r.screen, buttonX, screenHeight-2, question.Text, i == r.question.index) + 1
	}

	r.screen.Show()
}

func (r *Counter) renderRepoStatus(y int, repoStatus *repoStatus) {
	description := repoStatus.action.Description()
	if repoStatus.action == ActionError {
		if repoStatus.err == mgerrors.ErrNoChange {
			description = "no change"
		} else if repoStatus.err == mgerrors.ErrRejected {
			description = "rejected"
		} else if err, ok := repoStatus.err.(*exec.ExitError); ok {
			description = fmt.Sprintf("exit %d", err.ExitCode())
		}
	}

	emitStr(r.screen, r.columnStart(0), y, tcell.StyleDefault, shortenRepoName(repoStatus, maxRepoNameLen))
	emitStr(r.screen, r.columnStart(1), y, repoStatus.action.Color(), center(description, descriptionLen))
	progressBar(
		r.screen,
		20,
		repoStatus.action.Percentage(),
		r.columnStart(2),
		y,
	)
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
