package main

import "github.com/stewart/cmus"

// checks if two cmus status messages are different
func isDifferentStatus(a, b *cmus.Status) bool {
	if a.Playing != b.Playing {
		return true
	}

	if a.File != b.File {
		return true
	}

	if a.Duration != b.Duration || a.Position != b.Position {
		return true
	}

	if len(a.Tags) != len(b.Tags) {
		return false
	}

	for k, v := range a.Tags {
		if w, ok := b.Tags[k]; !ok || v != w {
			return false
		}
	}

	if len(a.Settings) != len(b.Settings) {
		return false
	}

	for k, v := range a.Settings {
		if w, ok := b.Settings[k]; !ok || v != w {
			return false
		}
	}

	return false
}
