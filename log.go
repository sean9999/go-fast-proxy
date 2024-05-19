package main

import (
	"cloud.google.com/go/logging"
)

func (d *Doggy) Slog(payload any, severity logging.Severity) {
	logger := d.Log.Logger("app")
	defer logger.Flush() // Ensure the entry is written.

	logger.Log(logging.Entry{
		// Log anything that can be marshaled to JSON.
		Payload:  payload,
		Severity: severity,
	})
}
