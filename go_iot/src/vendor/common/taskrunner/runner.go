package taskrunner

import (
	"fmt"
)

type Runner struct {
	Name       string
	Controller ControlChan
	Error      ControlChan
	Data       DataChan
	dataSize   int
	longlived  bool
	Dispatcher Fn
	Executor   Fn
}

func NewRunner(name string, size int, longlived bool, d Fn, e Fn) *Runner {
	return &Runner{
		Name:       name,
		Controller: make(chan string, 1),
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		dataSize:   size,
		longlived:  longlived,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longlived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	err := r.Dispatcher(r.Data)
	if err != nil {
		fmt.Printf("%s Error occur when dispatch data \n", r.Name)
	}
}

func (r *Runner) startAll() {
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}
