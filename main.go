package fastproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "google.golang.org/grpc/balancer/rls"
	_ "google.golang.org/grpc/xds/googledirectpath"
)

const bucketName = "go-proxy-02"

func init() {
	functions.HTTP("HelloHTTP", helloHTTP)
}

func helloHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	var d struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Hello, World!")
		return
	}
	if d.Name == "" {
		fmt.Fprint(w, "Hello, World!")
		return
	}

	fileName := fmt.Sprintf("%s.txt", d.Name)
	handle := client.Bucket(bucketName).Object(fileName)
	fw := handle.NewWriter(ctx)

	msg := []byte(fmt.Sprintf("my name is %s\n", d.Name))
	i, err := fw.Write(msg)
	defer fw.Close()

	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		fmt.Fprintf(w, "%d bytes were written to %s", i, fileName)
	}

}
