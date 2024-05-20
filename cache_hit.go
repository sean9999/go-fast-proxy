package main

import (
	"io"
	"net/http"

	logging "cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

func cacheHit(requestUri string, key string, d *Doggy, cacheReader *storage.Reader, httpWriter http.ResponseWriter) {

	//	object exists. Read from cache
	defer cacheReader.Close()
	merr := map[string]any{
		"attrs": cacheReader.Attrs,
		"msg":   "cache hit",
		"key":   key,
		"req":   requestUri,
	}
	d.Slog(merr, logging.Debug)
	io.Copy(httpWriter, cacheReader)

}
