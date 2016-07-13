package main

import (
	"sync"

	"github.com/stewart/cmus"
)

type State struct {
	sync.RWMutex

	status *cmus.Status
	err    error
}
