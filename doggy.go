package main

import (
	"context"
	"log"

	logging "cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

type Doggy struct {
	Ctx    context.Context
	Log    *logging.Client
	Store  *storage.Client
	Pubsub *pubsub.Client
}

func (d *Doggy) Teardown() {
	d.Log.Close()
	d.Store.Close()
	d.Pubsub.Close()
}

func NewDoggy(ctx context.Context) *Doggy {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	loggingClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		storageClient.Close()
		log.Fatal(err)
	}

	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		storageClient.Close()
		log.Fatal(err)
	}

	d := &Doggy{
		Ctx:    ctx,
		Store:  storageClient,
		Log:    loggingClient,
		Pubsub: pubSubClient,
	}
	return d
}
