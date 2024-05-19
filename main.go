package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

const projectID = "proxy02-423811"
const BUCKET = "go-proxy-cache-hash"

var loggingClient *logging.Client
var storageClient *storage.Client
var ctx context.Context

func init() {
	ctx = context.Background()
	storageClient, _ = storage.NewClient(ctx)
	loggingClient, _ = logging.NewClient(ctx, projectID)
}

func main() {

	log.Print("starting server...")
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

func handler(httpWriter http.ResponseWriter, httpReader *http.Request) {
	requestUri := httpReader.URL.RequestURI()

	ctx := context.Background()

	m5 := md5.New()
	io.WriteString(m5, requestUri)
	fmt.Printf("hash is %x", m5.Sum(nil))
	key := fmt.Sprintf("md5/%x", m5.Sum(nil))

	rc, err := storageClient.Bucket(BUCKET).Object(key).NewReader(ctx)
	if err != nil {

		//	object doesn't exist. Fetch and write

		merr := map[string]any{
			"error":   err,
			"context": "we tried to open a reader but it failed",
		}

		Slog(ctx, merr, logging.Info)

		bucketWriter := storageClient.Bucket(BUCKET).Object(key).NewWriter(ctx)
		defer bucketWriter.Close()

		client := &http.Client{}

		newAddress := fmt.Sprintf("https://goproxy.io%s", httpReader.RequestURI)

		log.Printf("newAddress is %s", newAddress)

		redir, err := http.NewRequestWithContext(ctx, http.MethodGet, newAddress, nil)
		if err != nil {

			merr := map[string]any{
				"error":    err,
				"contenxt": "we tried to create a new request object",
			}

			Slog(ctx, merr, logging.Alert)
			log.Fatal(err)
		}

		resp, err := client.Do(redir)
		if err != nil {

			merr := map[string]any{
				"error":    err,
				"contenxt": "httpClient failed to Do()",
			}

			Slog(ctx, merr, logging.Alert)
			log.Fatal(err)
		}

		//	write to both bucket and http response
		r2 := io.TeeReader(resp.Body, bucketWriter)
		i, err := io.Copy(httpWriter, r2)

		merr = map[string]any{
			"bytes_written": i,
			"contenxt":      "operation seems successful. We wrote the the bucket and to the http response",
		}
		Slog(ctx, merr, logging.Info)

		log.Printf("%d bytes written", i)
		if err != nil {
			log.Fatal(err)
		}

	} else {

		//	object exists. Read from cache
		log.Println("CACHE HIT!")
		merr := map[string]any{
			"msg":      "CACHE HIT",
			"contenxt": "we read from Google Storage",
		}
		Slog(ctx, merr, logging.Info)

		defer rc.Close()
		io.Copy(httpWriter, rc)
	}

	merr := map[string]any{
		"msg":        "lifecycle complete",
		"requestUri": requestUri,
		"key":        key,
	}
	Slog(ctx, merr, logging.Debug)

	log.Printf("The requestUri was %s and the hash is %s\n", requestUri, key)
}
