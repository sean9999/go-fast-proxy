# GCP Go Proxy

This app is meant to be deployed to Cloud Run on GCP. It saves it's data to Cloud Storage. You will need to set up a bucket.

The following environment variables must be set:

- PROJECT_ID
- STORAGE_BUCKET

Once deployed, you use the application URL in `GOPROXY`.

