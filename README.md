# GCP Go Proxy

This app is meant to be deployed to Cloud Run on GCP. It saves it's data to Cloud Storage. You will need to set up a bucket.

It first checks it's cache in the Storage Bucket. If it's there (cache hit), it simply returns the data.

If no entry is found (cache miss), it asks an upstream go proxy server, caches (save to bucket), and returns the result.

Caches stay forever, you'll want to configure a lifecycle policy on your Cloud Storage Bucket.

The following environment variables must be set:

- PROJECT_ID
- STORAGE_BUCKET

You may also set these, if you don't want the default values:

- PORT
- UPSTREAM_SERVER

Once deployed, you use the application URL in `GOPROXY`. Ex:

```sh
$ export GOPROXY='https://fast-go-proxy-qrpdktqpkq-nn.a.run.app,direct'
```
