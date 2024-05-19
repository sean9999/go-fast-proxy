package main

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
)

func Slog(ctx context.Context, payload any, severity logging.Severity) {
	if loggingClient == nil {
		log.Fatal("Failed to create logging client")
	}
	defer loggingClient.Close()
	logger := loggingClient.Logger("log")
	defer logger.Flush() // Ensure the entry is written.

	logger.Log(logging.Entry{
		// Log anything that can be marshaled to JSON.
		Payload:  payload,
		Severity: severity,
	})
}
