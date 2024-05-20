package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"

	"cloud.google.com/go/firestore"
	logging "cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

type Doggy struct {
	Ctx   context.Context
	Log   *logging.Client
	Store *storage.Client
	Fire  *firestore.Client
}

func (d *Doggy) Teardown() {
	d.Log.Close()
	d.Store.Close()
	d.Fire.Close()
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

	conf := &firebase.Config{ProjectID: projectID, DatabaseURL: fireBaseDb}
	app, _ := firebase.NewApp(ctx, conf)
	// if err != nil {
	// 	storageClient.Close()
	// 	loggingClient.Close()
	// 	log.Fatal(err)
	// }

	fire, _ := app.Firestore(ctx)
	// if err != nil {
	// 	storageClient.Close()
	// 	loggingClient.Close()
	// 	log.Fatal(err)
	// }

	// fire, err := firestore.NewClientWithDatabase(ctx, projectID, fireBaseDb)
	// if err != nil {
	// 	storageClient.Close()
	// 	loggingClient.Close()
	// 	log.Fatal(err)
	// }

	// fire, err := app.Firestore(ctx)
	// if err != nil {
	// 	storageClient.Close()
	// 	loggingClient.Close()
	// 	log.Fatal(err)
	// }

	d := &Doggy{
		Ctx:   ctx,
		Store: storageClient,
		Log:   loggingClient,
		Fire:  fire,
	}
	return d
}
