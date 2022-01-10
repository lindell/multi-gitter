package repocounter

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Action int

const (
	ActionNotStarted Action = iota
	ActionClone
	ActionRun
	ActionPush
	ActionCreatePR

	ActionSuccess
	ActionError
)

func (a Action) Description() string {
	switch a {
	case ActionNotStarted:
		return "waiting"
	case ActionClone:
		return "cloning"
	case ActionRun:
		return "running"
	case ActionPush:
		return "pushing"
	case ActionCreatePR:
		return "creating pr"
	case ActionSuccess:
		return "success"
	case ActionError:
		return "error"
	default:
		panic(fmt.Sprintf("could not get description of action %q", a))
	}
}

func (a Action) Color() tcell.Style {
	switch a {
	case ActionClone:
		return tcell.StyleDefault.Background(tcell.ColorBlue)
	case ActionRun:
		return tcell.StyleDefault.Background(tcell.ColorBlue)
	case ActionPush:
		return tcell.StyleDefault.Background(tcell.ColorBlue)
	case ActionCreatePR:
		return tcell.StyleDefault.Background(tcell.ColorBlue)
	case ActionSuccess:
		return tcell.StyleDefault.Background(tcell.ColorForestGreen)
	case ActionError:
		return tcell.StyleDefault.Background(tcell.ColorRed)
	default:
		return tcell.StyleDefault.Background(tcell.ColorDefault)
	}
}

func statusToPercentage(s Action) float64 {
	switch s {
	case ActionNotStarted:
		return 0.0
	case ActionClone:
		return 0.1
	case ActionRun:
		return 0.3
	case ActionPush:
		return 0.6
	case ActionCreatePR:
		return 0.9
	case ActionSuccess:
		return 1.0
	case ActionError:
		return 1.0
	default:
		panic(fmt.Sprintf("could not get percentage of status %q", s))
	}
}
