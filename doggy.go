package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	logging "cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

type Doggy struct {
	Ctx     context.Context
	Logging *logging.Client
	Storing *storage.Client
	Burning *firestore.Client
}

func (d *Doggy) Teardown() {
	d.Logging.Close()
	d.Storing.Close()
	d.Burning.Close()
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

	fireClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		storageClient.Close()
		loggingClient.Close()
		log.Fatal(err)
	}
	d := &Doggy{
		Ctx:     ctx,
		Storing: storageClient,
		Logging: loggingClient,
		Burning: fireClient,
	}
	return d
}
