package main

import (
	"errors"
	"fmt"
	"os"
)

var projectID = os.Getenv("PROJECT_ID")
var storageBucket = os.Getenv("STORAGE_BUCKET")
var port = os.Getenv("PORT")

const upstreamServer = "https://goproxy.io"

var ErrBadSettings = errors.New("bad settings")

func defaults() error {
	if projectID == "" {
		return fmt.Errorf("%w. PROJECT_ID environment variable needs to be set", ErrBadSettings)
	}
	if storageBucket == "" {
		return fmt.Errorf("%w. STORAGE_BUCKET needs to be set", ErrBadSettings)
	}
	if port == "" {
		port = "8080"
	}
	return nil
}
