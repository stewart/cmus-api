package main

import "github.com/stewart/cmus"

// given two cmus statuses, generates a map of the difference between them
func diffStatus(a, b *cmus.Status) map[string]interface{} {
	diff := make(map[string]interface{})

	var hashDiff = func(x, y map[string]string) map[string]string {
		z := make(map[string]string)

		for k, v := range x {
			w, ok := y[k]

			if !ok {
				z[k] = ""
			}

			if v != w {
				z[k] = w
			}
		}

		for k, v := range y {
			w, ok := x[k]

			if !ok {
				z[k] = ""
			}

			if v != w {
				z[k] = v
			}
		}

		return z
	}

	if a.Playing != b.Playing {
		diff["playing"] = b.Playing
	}

	if a.File != b.File {
		diff["file"] = b.File
	}

	if a.Duration != b.Duration {
		diff["duration"] = b.Duration
	}

	if a.Position != b.Position {
		diff["position"] = b.Position
	}

	tags := hashDiff(a.Tags, b.Tags)
	settings := hashDiff(a.Settings, b.Settings)

	if len(tags) > 0 {
		diff["tags"] = tags
	}

	if len(settings) > 0 {
		diff["settings"] = settings
	}

	return diff
}

// converts a status into a map that Gin can serialize
func serializeStatus(s *cmus.Status) map[string]interface{} {
	return map[string]interface{}{
		"playing":  s.Playing,
		"file":     s.File,
		"duration": s.Duration,
		"position": s.Position,
		"tags":     s.Tags,
		"settings": s.Settings,
	}
}
