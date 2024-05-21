package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
)

func cacheMiss(requestUri string, key string, d *Doggy, httpReader *http.Request, httpWriter http.ResponseWriter) {

	//	queue the job
	topic, err := d.Pubsub.Topic(d.Ctx, pubsubTopic)
	if err != nil {
		log.Fatal(err)
	}
	payload := map[string]string{
		"requestUri": requestUri,
		"key":        key,
	}
	j, _ := json.Marshal(payload)
	topic.Publish(d.Ctx, &pubsub.Message{
		Data: []byte(j),
	})

	//	redirect to upstream
	newAddress := fmt.Sprintf("%s%s", upstreamServer, httpReader.RequestURI)
	httpWriter.Header().Set("Content-Type", "")
	http.Redirect(httpWriter, httpReader, newAddress, http.StatusTemporaryRedirect)

}
