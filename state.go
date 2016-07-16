package main

import (
	"sync"

	"github.com/stewart/cmus"
)

type State struct {
	sync.RWMutex
	cmus   *cmus.Client
	status *cmus.Status
	err    error
}

var state = &State{}
