package main

import (
	"context"
	"log"

	logging "cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

// Workhorse is our singleton that contains references to everything we need
type Workhorse struct {
	Ctx   context.Context
	Log   *logging.Client
	Store *storage.Client
}

// This should be invoked at the end of the lifecycle
func (app *Workhorse) Teardown() {
	app.Log.Close()
	app.Store.Close()
}

// create a new App. Die if anything goes wrong
func NewWorkhorse(ctx context.Context) *Workhorse {

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	loggingClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		storageClient.Close()
		log.Fatal(err)
	}

	app := &Workhorse{
		Ctx:   ctx,
		Store: storageClient,
		Log:   loggingClient,
	}
	return app
}
