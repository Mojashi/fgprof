package fgprof

import (
	"bytes"
	"net/http"
)

// Handler returns an http handler that takes an optional "seconds" query
// argument that defaults to "30" and produces a profile over this duration.
// The optional "format" parameter controls if the output is written in
// Google's "pprof" format (default) or Brendan Gregg's "folded" stack format.

type ProfileInstance struct {
	stopper func() error
	writer  *bytes.Buffer
}

var runningProfile *ProfileInstance = nil

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		format := FormatPprof

		switch r.URL.Query().Get("command") {
		case "start":
			if runningProfile != nil {
				runningProfile.stopper()
			}

			writer := &bytes.Buffer{}
			stopper := Start(writer, format)
			runningProfile = &ProfileInstance{stopper, writer}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("profile started\n"))

		case "stop":
			if runningProfile == nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("no profile running\n"))
				return
			}
			runningProfile.stopper()
			w.WriteHeader(http.StatusOK)
			runningProfile.writer.WriteTo(w)
			runningProfile = nil

		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("unknown command\n"))
		}
	})
}
