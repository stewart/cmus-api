package main

import (
	"time"

	"github.com/stewart/cmus"
)

func poll() (<-chan (map[string]interface{}), <-chan error) {
	diffs := make(chan map[string]interface{})
	errs := make(chan error)

	go func() {
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
					diffs <- serializeStatus(status)
					initial = false
				} else if prevErr != nil || len(diff) > 0 {
					diffs <- diff
				}
			} else {
				// if a new error, send error
				if prevErr == nil || prevErr.Error() != err.Error() {
					errs <- err
				}
			}

			prevStatus = status
			prevErr = err

			state.RUnlock()

			time.Sleep(50 * time.Millisecond)
		}
	}()

	return diffs, errs
}
