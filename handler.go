package main

import (
	"errors"
	"net/http"
	"path"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/storage"
)

func (d *Workhorse) ServeHTTP(httpWriter http.ResponseWriter, httpReader *http.Request) {

	requestUri := httpReader.URL.RequestURI()
	key := path.Join("plain", requestUri)

	obj := d.Store.Bucket(storageBucket).Object(key)
	cacheReader, err := obj.NewReader(d.Ctx)

	if errors.Is(storage.ErrObjectNotExist, err) {

		cacheMiss(requestUri, obj, key, d, httpReader, httpWriter)

	} else if err != nil {

		merr := map[string]any{
			"err": err,
			"msg": "there was an error reading the object from storage, but it wasn't ErrObjectNotExist",
		}
		d.Slog(merr, logging.Alert)

		floo{
			"nerd": true,
			"age":  48,
			"ok":   nil,
		}.Validate().Log()

	} else {

		cacheHit(requestUri, key, d, cacheReader, httpWriter)

	}

}
