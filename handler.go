package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"

	logging "cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

func (d *Doggy) ServeHTTP(httpWriter http.ResponseWriter, httpReader *http.Request) {

	requestUri := httpReader.URL.RequestURI()
	key := path.Join("plain", requestUri)

	o := d.Store.Bucket(storageBucket).Object(key)
	rc, err := o.NewReader(d.Ctx)

	if err != nil {

		merr := map[string]any{
			"msg": "cache miss",
			"key": key,
		}
		d.Slog(merr, logging.Debug)

		//	create a bucket writer

		bucketWriter := o.NewWriter(d.Ctx)
		bucketWriter.ObjectAttrs = storage.ObjectAttrs{Metadata: map[string]string{
			"requestUri": requestUri,
			"key":        key,
		}}

		//	create a new HTTP request to upstream server
		client := &http.Client{}
		newAddress := fmt.Sprintf("https://goproxy.io%s", httpReader.RequestURI)
		log.Printf("newAddress is %s", newAddress)
		redir, err := http.NewRequestWithContext(d.Ctx, http.MethodGet, newAddress, nil)
		if err != nil {
			merr := map[string]any{
				"error":      err,
				"msg":        "we tried to create a new request object",
				"key":        key,
				"requestUri": requestUri,
			}
			d.Slog(merr, logging.Alert)
			log.Fatal(err)
		}

		//	issue the upstream request
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

		//	pipe the response to our upstream request, to bucketWriter _and_ the main http.Response

		r2 := io.TeeReader(resp.Body, bucketWriter)

		i, err := io.Copy(httpWriter, r2)
		if err != nil {
			log.Fatalln(err)
		}

		defer bucketWriter.Close()

		// merr = map[string]any{
		// 	"bytes_written": i,
		// 	"msg":           "operation seems successful. We wrote the the bucket and to the http response",
		// 	"m5str":         m5str,
		// 	"key":           key,
		// 	"requestUri":    requestUri,
		// }
		// d.Slog(merr, logging.Info)
		log.Printf("%d bytes written", i)
		if err != nil {
			log.Fatal(err)
		}

		//	if all is good, write to FireStore
		// caches := d.Burning.Collection("goproxy-cache-lookup/caches")
		// thisDoc := caches.Doc(key)
		// thisDoc.Create(d.Ctx, CacheTuple{
		// 	Req:   requestUri,
		// 	Hash:  key,
		// 	Mtime: uint64(time.Now().UnixMicro()),
		// 	Atime: uint64(time.Now().UnixMicro()),
		// })

	} else {

		//	object exists. Read from cache
		defer rc.Close()

		log.Println("CACHE HIT!")
		// merr := map[string]any{
		// 	"msg":        "CACHE HIT",
		// 	"contenxt":   "we read from Google Storage",
		// 	"requestUri": requestUri,
		// }
		// d.Slog(merr, logging.Info)

		io.Copy(httpWriter, rc)

	}

	// merr := map[string]any{
	// 	"msg":        "lifecycle complete",
	// 	"requestUri": requestUri,
	// 	"key":        key,
	// }
	// d.Slog(merr, logging.Debug)

	log.Printf("The requestUri was %s and the hash is %s\n", requestUri, key)
}
