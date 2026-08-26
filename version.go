package main

import (
	"fmt"
	"io"
	"runtime/debug"
)

func printVersion(stdout io.Writer) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintln(stdout, "dl version unknown")
		return
	}

	var tag, hash string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.tag":
			tag = s.Value
		case "vcs.revision":
			hash = s.Value
		}
	}

	switch {
	case tag != "" && hash != "":
		shortHash := hash
		if len(hash) > 7 {
			shortHash = hash[:7]
		}
		fmt.Fprintf(stdout, "dl %s (%s)\n", tag, shortHash)
	case tag != "":
		fmt.Fprintf(stdout, "dl %s\n", tag)
	case hash != "":
		shortHash := hash
		if len(hash) > 7 {
			shortHash = hash[:7]
		}
		fmt.Fprintf(stdout, "dl %s\n", shortHash)
	default:
		fmt.Fprintln(stdout, "dl version unknown")
	}
}
