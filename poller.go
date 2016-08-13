package main

import (
	"time"

	"github.com/stewart/cmus"
)

type Poller struct {
	Diffs  chan map[string]interface{}
	Errors chan error

	done chan int
}

func NewPoller() *Poller {
	p := &Poller{
		Diffs:  make(chan map[string]interface{}),
		Errors: make(chan error),
		done:   make(chan int),
	}

	return p
}

func (p *Poller) Poll() {
	initial := true
	prevStatus := &cmus.Status{}
	var prevErr error

	for {
		state.RLock()

		status := state.status
		err := state.err

		if err == nil {
			diff := diffStatus(prevStatus, status)

			if initial {
				select {
				case p.Diffs <- serializeStatus(status):
				case <-p.done:
					state.RUnlock()
					return
				}
				initial = false
			} else if prevErr != nil || len(diff) > 0 {
				select {
				case p.Diffs <- diff:
				case <-p.done:
					state.RUnlock()
					return
				}
			}
		} else {
			// if a new error, send error
			if prevErr == nil || prevErr.Error() != err.Error() {
				select {
				case p.Errors <- err:
				case <-p.done:
					state.RUnlock()
					return
				}
			}
		}

		prevStatus = status
		prevErr = err

		state.RUnlock()

		time.Sleep(50 * time.Millisecond)
	}
}

func (p *Poller) Close() {
	close(p.done)
}
