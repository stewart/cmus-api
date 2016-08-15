package main

import (
	"time"

	"github.com/stewart/cmus"
)

type Poller struct {
	Diffs  chan map[string]interface{}
	Errors chan error

	initialized    bool
	previousStatus *cmus.Status
	previousError  error

	tick *time.Ticker
	done chan int
}

func NewPoller() *Poller {
	p := &Poller{
		Diffs:  make(chan map[string]interface{}),
		Errors: make(chan error),

		previousStatus: &cmus.Status{},
		done:           make(chan int),
	}

	return p
}

func (p *Poller) Poll() {
	p.tick = time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-p.tick.C:
			p.update()
		case <-p.done:
			p.tick.Stop()
			return
		}
	}
}

func (p *Poller) update() {
	state.RLock()
	defer state.RUnlock()

	status := state.status
	err := state.err

	if err == nil {
		diff := diffStatus(p.previousStatus, status)

		if !p.initialized {
			p.Diffs <- serializeStatus(status)
			p.initialized = true
		} else if p.previousError != nil || len(diff) > 0 {
			p.Diffs <- diff
		}
	} else {
		// if a new error, send error
		if p.previousError == nil || p.previousError.Error() != err.Error() {
			p.Errors <- err
		}
	}

	p.previousStatus = status
	p.previousError = err
}

func (p *Poller) Close() {
	close(p.done)
}
