package ui

import (
	"github.com/cinarmert/klogs/cmd/klogs/pod"
	"github.com/gdamore/tcell"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"os"
	"sync"
)

// LayoutManager sets the layout according to given targets
type LayoutManager struct {
	App     *tview.Application
	Grid    *tview.Grid
	Targets []*pod.Target
	core    v1.CoreV1Interface
}

// NewManager return a manager for the given Targets and k8s CoreV1 client.
func NewManager(sessions []*pod.Target, cs v1.CoreV1Interface) (*LayoutManager, error) {
	if len(sessions) == 0 {
		return nil, errors.New("could not create layout manager: no targets given")
	}

	lm := (&LayoutManager{}).
		setApp().
		setSessions(sessions).
		setGrid().
		setClientSet(cs)

	return lm, nil
}

func (lm *LayoutManager) setApp() *LayoutManager {
	lm.App = tview.NewApplication()
	return lm
}

func (lm *LayoutManager) setClientSet(core v1.CoreV1Interface) *LayoutManager {
	lm.core = core
	return lm
}

func (lm *LayoutManager) setSessions(sessions []*pod.Target) *LayoutManager {
	lm.Targets = sessions
	return lm
}

func (lm *LayoutManager) setGrid() *LayoutManager {
	n := len(lm.Targets)
	rows := (n + 1) / 2
	cols := 2

	if n == 1 {
		cols = 1
	}

	grid := tview.NewGrid().
		SetRows(make([]int, rows)...).
		SetColumns(make([]int, cols)...).
		SetBorders(true)

	lm.Grid = grid
	return lm
}

func (lm *LayoutManager) createTextView(title string) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		ScrollToEnd().
		SetChangedFunc(func() {
			lm.App.Draw()
		})

	textView.SetBorder(true)
	textView.SetTitle(" " + title + " ")
	textView.SetTitleColor(tcell.ColorRed)
	return textView
}

// Run starts the log and ui threads.
func (lm *LayoutManager) Run() {
	var wg sync.WaitGroup
	for i, target := range lm.Targets {
		wg.Add(1)
		tv := lm.createTextView(target.Pod + "/" + target.Container)

		row, col, colSpan := i/2, i%2, 1
		if i == len(lm.Targets)-1 && len(lm.Targets)%2 == 1 {
			colSpan = 2
		}

		lm.Grid.AddItem(tv, row, col, 1, colSpan, 0, 0, false)

		go target.StartThread(lm.core, &wg, tview.ANSIWriter(tv))
	}

	if !isatty.IsTerminal(os.Stdout.Fd()) {
		log.Warnf("klogs is only available in a tty environment")
		return
	}

	if err := lm.App.SetRoot(lm.Grid, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("could not init ui: %v", err)
		os.Exit(1)
	}
}
