package main

import "time"

// polls cmus for status, updating state as appropriate
func loop() {
	go func() {
		for {
			if state.err == nil {
				fetch()
			} else {
				reconnect()
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()
}

func fetch() {
	status, err := state.cmus.Status()

	state.Lock()
	state.status = status
	state.err = err
	state.Unlock()
}

func reconnect() {
	state.Lock()
	state.Unlock()
	state.err = state.cmus.Connect()
}
