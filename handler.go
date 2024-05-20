package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"

	logging "cloud.google.com/go/logging"
)

func (d *Doggy) ServeHTTP(httpWriter http.ResponseWriter, httpReader *http.Request) {

	requestUri := httpReader.URL.RequestURI()
	key := path.Join("plain", requestUri)

	o := d.Store.Bucket(storageBucket).Object(key)
	cacheReader, err := o.NewReader(d.Ctx)

	if err != nil {

		//	create a bucket writer
		bucketWriter := o.NewWriter(d.Ctx)

		//	create a new HTTP request to upstream server
		client := &http.Client{}
		newAddress := fmt.Sprintf("https://goproxy.io%s", httpReader.RequestURI)

		redir, err := http.NewRequestWithContext(d.Ctx, http.MethodGet, newAddress, nil)
		if err != nil {
			merr := map[string]any{
				"error":      err,
				"msg":        "we tried to create a new request object, but it failed",
				"key":        key,
				"requestUri": requestUri,
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

		//	pipe the response to our upstream request
		//	to bucketWriter and the main http.Response
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

	} else {

		//	object exists. Read from cache
		defer cacheReader.Close()
		merr := map[string]any{
			"attrs": cacheReader.Attrs,
			"msg":   "cache hit",
		}
		d.Slog(merr, logging.Debug)
		io.Copy(httpWriter, cacheReader)

	}

}
