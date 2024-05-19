package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
)

type FastEvent struct {
	logging.Entry
	Payload string
}

const projectID = "projects/proxy02-423811"

func main() {

	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.

	// Creates a logClient.
	logClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer logClient.Close()

	logger := logClient.Logger("hello")

	logger.Log(logging.Entry{Payload: "now we have a logger"})

	log.Print("starting server...")
	http.HandleFunc("/", handler)

	logger.Log(logging.Entry{Severity: logging.Info, Payload: "starting server"})

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
		logger.Log(logging.Entry{Severity: logging.Alert, Payload: "defaulting to port 8080"})
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	logger.Log(logging.Entry{Severity: logging.Error, Payload: "end of run loop"})

}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}
