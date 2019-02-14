package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type Runner struct {
	interrupt chan os.Signal
	complete  chan error
	timeout   <-chan time.Time
	tasks     []func(int)
}

//ErrTimeout 任务执行超时的时候返回
var ErrTimeout = errors.New("received timeout")

//ErrInterrupt 会在接收到操作系统的事件时返回
var ErrInterrupt = errors.New("received interrupt")

func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() error {
	//接收终端信号
	signal.Notify(r.interrupt, os.Interrupt)

	go func() {
		r.complete <- r.run()
	}()

	go func() {
		r.interrupt <- os.Kill
	}()

	select {
	case err := <-r.complete:
		return err
	case <-r.timeout:
		return ErrInterrupt
	}
}

func (r *Runner) run() error {
	for id, task := range r.tasks {
		if r.gotInterrupt() {
			return ErrInterrupt
		}
		task(id)
	}

	return nil
}

func (r *Runner) gotInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}
