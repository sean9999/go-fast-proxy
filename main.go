package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
)

const BUCKET = "go-proxy-cache-hash"

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

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Read the object1 from bucket.
	rc, err := client.Bucket(BUCKET).Object(key).NewReader(ctx)
	if err != nil {

		log.Println(err)

		bucketWriter := client.Bucket(BUCKET).Object(key).NewWriter(ctx)

		client := &http.Client{}
		redir := httpReader.WithContext(ctx)
		redir.URL.Host = "goproxy.io"
		redir.Host = "goproxy.io"

		resp, err := client.Do(redir)
		if err != nil {
			log.Fatal(err)
		}

		//	write to both bucket and http response
		r2 := io.TeeReader(resp.Body, bucketWriter)
		i, err := io.Copy(httpWriter, r2)
		log.Printf("%d bytes written", i)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		defer rc.Close()
		io.Copy(httpWriter, rc)
	}

	log.Printf("The requestUri was %s and the hash is %s\n", requestUri, key)
	fmt.Fprintf(httpWriter, "The requestUri was %s!\n", requestUri)
}
