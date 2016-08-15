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
	for {
		state.RLock()

		status := state.status
		err := state.err

		if err == nil {
			diff := diffStatus(p.previousStatus, status)

			if !p.initialized {
				select {
				case p.Diffs <- serializeStatus(status):
				case <-p.done:
					state.RUnlock()
					return
				}
				p.initialized = true
			} else if p.previousError != nil || len(diff) > 0 {
				select {
				case p.Diffs <- diff:
				case <-p.done:
					state.RUnlock()
					return
				}
			}
		} else {
			// if a new error, send error
			if p.previousError == nil || p.previousError.Error() != err.Error() {
				select {
				case p.Errors <- err:
				case <-p.done:
					state.RUnlock()
					return
				}
			}
		}

		p.previousStatus = status
		p.previousError = err

		state.RUnlock()

		time.Sleep(50 * time.Millisecond)
	}
}

func (p *Poller) Close() {
	close(p.done)
}

func (p *Poller) update() {
}
