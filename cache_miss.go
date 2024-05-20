package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	logging "cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

func cacheMiss(requestUri string, obj *storage.ObjectHandle, key string, d *Doggy, httpReader *http.Request, httpWriter http.ResponseWriter) {

	//	create a bucket writer
	bucketWriter := obj.NewWriter(d.Ctx)

	//	create a new HTTP request to upstream server
	client := &http.Client{}
	newAddress := fmt.Sprintf("%s%s", upstreamServer, httpReader.RequestURI)
	redir, err := http.NewRequestWithContext(d.Ctx, http.MethodGet, newAddress, nil)
	if err != nil {
		merr := map[string]any{
			"error":           err,
			"msg":             "we tried to create a request, but it failed",
			"key":             key,
			"requestUri":      requestUri,
			"upstreamAddress": newAddress,
		}
		d.Slog(merr, logging.Error)
		log.Fatal(err)
	}

	//	send the upstream request
	resp, err := client.Do(redir)
	if err != nil {
		merr := map[string]any{
			"error":      err,
			"msg":        "httpClient failed to Do()",
			"key":        key,
			"addr":       newAddress,
			"requestUri": requestUri,
		}
		d.Slog(merr, logging.Alert)
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//	Pipe the response of our upstream request
	//	to bucketWriter and the main http.Response, simultaneously.
	r2 := io.TeeReader(resp.Body, bucketWriter)
	i, err := io.Copy(httpWriter, r2)
	if err != nil {
		log.Fatalln(err)
	}
	defer bucketWriter.Close()

	if err != nil {
		log.Fatal(err)
	} else {
		merr := map[string]any{
			"msg":             "cache miss",
			"key":             key,
			"upstreamAddress": newAddress,
			"bytesWritten":    i,
		}
		d.Slog(merr, logging.Debug)
	}

}
