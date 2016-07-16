package main

import (
	"testing"

	"github.com/stewart/cmus"
)

func TestIsDifferentStatus(t *testing.T) {
	cases := []struct {
		A, B     *cmus.Status
		Expected bool
	}{
		{
			&cmus.Status{},
			&cmus.Status{},
			false,
		},
		{
			&cmus.Status{Playing: true},
			&cmus.Status{Playing: false},
			true,
		},
		{
			&cmus.Status{Duration: 108, Position: 100},
			&cmus.Status{Duration: 108, Position: 200},
			true,
		},
		{
			&cmus.Status{Tags: map[string]string{"hello": "world"}},
			&cmus.Status{Tags: map[string]string{"hello": "world"}},
			false,
		},

		{
			&cmus.Status{Tags: map[string]string{"hello": "yep"}},
			&cmus.Status{Tags: map[string]string{"hello": "world"}},
			false,
		},
	}

	for _, c := range cases {
		got := isDifferentStatus(c.A, c.B)

		if got != c.Expected {
			msg := "Expected isDifferentStatus(%v, %v) to be %b, got %b"
			t.Errorf(msg, c.A, c.B, c.Expected, got)
		}
	}
}
