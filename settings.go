package main

import (
	"errors"
	"fmt"
	"os"
)

const PubsubLiteRegion = "us-central1"

var pubsubTopic = os.Getenv("PUBSUB_TOPIC")
var pubsubSubscription = os.Getenv("PUBSUB_SUB")
var projectID = os.Getenv("PROJECT_ID")
var storageBucket = os.Getenv("STORAGE_BUCKET")
var port = os.Getenv("PORT")
var upstreamServer = os.Getenv("UPSTREAM_SERVER")
var ErrBadSettings = errors.New("bad settings")

func defaults() error {
	if projectID == "" {
		return fmt.Errorf("%w. PROJECT_ID environment variable needs to be set", ErrBadSettings)
	}
	if storageBucket == "" {
		return fmt.Errorf("%w. STORAGE_BUCKET needs to be set", ErrBadSettings)
	}
	if pubsubTopic == "" {
		return fmt.Errorf("%w. PUBSUB_TOPIC needs to be set", ErrBadSettings)
	}
	if pubsubSubscription == "" {
		return fmt.Errorf("%w. PUBSUB_SUB needs to be set", ErrBadSettings)
	}
	if port == "" {
		port = "8080"
	}
	if upstreamServer == "" {
		upstreamServer = "https://goproxy.io"
	}
	return nil
}
