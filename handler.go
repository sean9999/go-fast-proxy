package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

type CacheTuple struct {
	Req   string `json:"req,omitempty"`
	Hash  string `json:"hash,omitempty"`
	Mtime uint64 `json:"mtime,omitempty"`
	Atime uint64 `json:"atime,omitempty"`
}

func (d *Doggy) ServeHTTP(httpWriter http.ResponseWriter, httpReader *http.Request) {
	requestUri := httpReader.URL.RequestURI()

	//	base64
	baseBuf := new(bytes.Buffer)
	input := []byte(requestUri)
	encoder := base64.NewEncoder(base64.StdEncoding, baseBuf)
	encoder.Write(input)

	//	hex
	hex := fmt.Sprintf("%x", requestUri)

	//	md5
	m5 := md5.New()
	io.WriteString(m5, hex)
	m5str := fmt.Sprintf("md5/%x", m5.Sum(nil))
	//key := fmt.Sprintf("hex/%x", requestUri)

	log.Println(hex, m5str)

	rc, err := d.Store.Bucket(storageBucket).Object(hex).NewReader(d.Ctx)
	if err != nil {

		log.Println(rc, err)

		log.Println("CACHE miss :(")

		//	object doesn't exist. Fetch and write

		// merr := map[string]any{
		// 	"error":      err,
		// 	"msg":        "we tried to open a reader but it failed",
		// 	"todo":       "check for specific type of error indicating 404",
		// 	"m5str":      m5str,
		// 	"key":        key,
		// 	"requestUri": requestUri,
		// }
		// d.Slog(merr, logging.Info)

		//	create a bucket writer

		o := d.Store.Bucket(storageBucket).Object(hex)

		bucketWriter := o.NewWriter(d.Ctx)
		// bucketWriter.ObjectAttrs = storage.ObjectAttrs{Metadata: map[string]string{
		// 	"requestUri": requestUri,
		// 	"key":        hex,
		// 	"nerd":       "poo",
		// 	"m5str":      m5str,
		// 	"base64":     baseBuf.String(),
		// }}

		//	create a new HTTP request to upstream server
		client := &http.Client{}
		newAddress := fmt.Sprintf("https://goproxy.io%s", httpReader.RequestURI)
		log.Printf("newAddress is %s", newAddress)
		redir, err := http.NewRequestWithContext(d.Ctx, http.MethodGet, newAddress, nil)
		if err != nil {
			// merr := map[string]any{
			// 	"error":      err,
			// 	"msg":        "we tried to create a new request object",
			// 	"m5str":      m5str,
			// 	"requestUri": requestUri,
			// }
			// d.Slog(merr, logging.Alert)
			log.Fatal(err)
		}

		//	issue the upstream request
		resp, err := client.Do(redir)
		if err != nil {
			// merr := map[string]any{
			// 	"error":      err,
			// 	"msg":        "httpClient failed to Do()",
			// 	"m5str":      m5str,
			// 	"key":        key,
			// 	"addr":       newAddress,
			// 	"requestUri": requestUri,
			// }
			// d.Slog(merr, logging.Alert)
			log.Fatal(err)
		}

		//	pipe the response to our upstream request, to bucketWriter _and_ the main http.Response

		r2 := io.TeeReader(resp.Body, bucketWriter)

		o.Update(d.Ctx, storage.ObjectAttrsToUpdate{
			Metadata: map[string]string{
				"requestUri": requestUri,
				"key":        hex,
				"nerd":       "poo",
				"m5str":      m5str,
				"base64":     baseBuf.String(),
			},
		})

		// Update the object to set the metadata.
		objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
			Metadata: map[string]string{
				"hex":  hex,
				"base": baseBuf.String(),
				"req":  requestUri,
			},
		}

		if _, err := o.Update(d.Ctx, objectAttrsToUpdate); err != nil {
			log.Fatalf("ObjectHandle(%q).Update: %s", o.BucketName(), err)
		}

		defer bucketWriter.Close()
		i, err := io.Copy(httpWriter, r2)

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

		log.Println("CACHE HIT!")
		// merr := map[string]any{
		// 	"msg":        "CACHE HIT",
		// 	"contenxt":   "we read from Google Storage",
		// 	"requestUri": requestUri,
		// }
		// d.Slog(merr, logging.Info)

		io.Copy(httpWriter, rc)
		rc.Close()

	}

	// merr := map[string]any{
	// 	"msg":        "lifecycle complete",
	// 	"requestUri": requestUri,
	// 	"key":        key,
	// }
	// d.Slog(merr, logging.Debug)

	log.Printf("The requestUri was %s and the hash is %s\n", requestUri, hex)
}
