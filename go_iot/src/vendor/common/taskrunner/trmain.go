package taskrunner

import (
	"time"
)

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) StartDispatch() {
	for {
		select {
		case <-w.ticker.C:
			w.runner.startDispatch()
		}
	}
}

func (w *Worker) StartExecutor() {
	for {
		select {
		case record := <-w.runner.Data:
			w.runner.startExecutor(record)
		}
	}
}
