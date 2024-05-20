package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
)

func cacheMiss(requestUri string, key string, d *Doggy, httpReader *http.Request, httpWriter http.ResponseWriter) {

	//	queue the job
	topic := d.Pubsub.Topic(pubsubTopic)
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
