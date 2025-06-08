package vcs

import (
	"fmt"
	"runtime/debug"
)

func Version() string {
	var revision string
	var modified bool
	var time string

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range buildInfo.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				modified = s.Value == "true"
			case "vcs.time":
				time = s.Value
			}
		}
	}

	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}

	return fmt.Sprintf("%s-%s", time, revision)
}