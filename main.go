package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/logging"
)

func main() {

	//	load settings and defaults from env vars
	if err := defaults(); err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.Lshortfile)
	ctx := context.Background()

	d := NewWorkhorse(ctx)
	defer d.Teardown()

	http.Handle("/", d)

	if port == "" {
		port = "8080"
		d.Slog(map[string]any{
			"port": port,
			"msg":  "read port for env",
		}, logging.Debug)
		log.Printf("defaulting to port %s", port)
	}

	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		merr := map[string]any{
			"err": err,
			"msg": "could not start http server",
		}
		d.Slog(merr, logging.Critical)
		log.Fatal(err)
	}

	//topic := d.Pubsub.Topic(pubsubTopic)

	//sub := d.Pubsub.Subscription("cache-queue-lite")

}
