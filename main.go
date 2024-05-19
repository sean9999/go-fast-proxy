package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

const projectID = "proxy02-423811"
const BUCKET = "go-proxy-cache-hash"

func main() {

	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	loggingClient, err := logging.NewClient(ctx, projectID)
	defer loggingClient.Close()
	if err != nil {
		log.Fatal(err)
	}
	d := &Doggy{
		Ctx:           ctx,
		StorageClient: storageClient,
		LoggingClient: loggingClient,
	}

	log.Print("starting server...")
	http.Handle("/", d)
	//http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
