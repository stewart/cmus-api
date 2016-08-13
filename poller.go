package main

import (
	"time"

	"github.com/stewart/cmus"
)

type Poller struct {
	Diffs  chan map[string]interface{}
	Errors chan error
}

func NewPoller() *Poller {
	p := &Poller{
		Diffs:  make(chan map[string]interface{}),
		Errors: make(chan error),
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
				p.Diffs <- serializeStatus(status)
				initial = false
			} else if prevErr != nil || len(diff) > 0 {
				p.Diffs <- diff
			}
		} else {
			// if a new error, send error
			if prevErr == nil || prevErr.Error() != err.Error() {
				p.Errors <- err
			}
		}

		prevStatus = status
		prevErr = err

		state.RUnlock()

		time.Sleep(50 * time.Millisecond)
	}
}
