package main

import (
	"log"
	"time"
)

func loop() {
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	for {
		status, err := client.Status()

		state.Lock()

		if err == nil {
			if isDifferentStatus(state.status, status) {
				state.status = status
			}
		} else {
			// error when fetching status, attempt to reconnect
			client.Connect()
		}

		state.err = err

		state.Unlock()

		time.Sleep(100 * time.Millisecond)
	}
}
