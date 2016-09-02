package main

import (
	"time"

	"github.com/stewart/cmus"
)

type Poller struct {
	Statuses chan map[string]interface{}
	Errors   chan error

	IncrementalUpdates bool

	initialized    bool
	previousStatus *cmus.Status
	previousError  error

	tick *time.Ticker
	done chan int
}

// Creates a new Poller.
// `incrementalUpdates` determines whether or not Statuses will be sent
// incremental status diffs, or the full status on each change.
func NewPoller(incrementalUpdates bool) *Poller {
	p := &Poller{
		Statuses: make(chan map[string]interface{}),
		Errors:   make(chan error),

		IncrementalUpdates: incrementalUpdates,

		previousStatus: &cmus.Status{},
		done:           make(chan int),
	}

	return p
}

func (p *Poller) Poll() {
	p.tick = time.NewTicker(200 * time.Millisecond)

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
			p.Statuses <- serializeStatus(status)
			p.initialized = true
		} else if p.previousError != nil || len(diff) > 0 {
			if p.IncrementalUpdates {
				p.Statuses <- diff
			} else {
				p.Statuses <- serializeStatus(status)
			}
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
