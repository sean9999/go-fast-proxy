package main

import (
	"fmt"

	"cloud.google.com/go/logging"
)

func (d *Workhorse) Slog(payload any, severity logging.Severity) {
	logger := d.Log.Logger("app")
	defer logger.Flush() // Ensure the entry is written.

	logger.Log(logging.Entry{
		// Log anything that can be marshaled to JSON.
		Payload:  payload,
		Severity: severity,
	})
}

type floo map[string]any

func (f floo) Validate() floo {
	_, keyexists := f["nerd"]
	if !keyexists {
		panic("you must declare whether or not you're a nerd")
	}
	return f
}

func (f floo) Log() {
	fmt.Printf("%#v,\n", f)
}
