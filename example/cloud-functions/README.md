# Example

## How to run with Cloud Functions/Cloud Scheduler

```sh
$ gcloud pubsub topics create bqnotify

$ gcloud functions deploy BqNotify \
    --region asia-northeast1 \
    --runtime go113 \
    --trigger-topic bqnotify \
    --set-env-vars SLACK_WEBHOOK_URL=$SLACK_WEBHOOK_URL \
    --set-env-vars GCP_PROJECT=mizzy-270104 # GCP_PRJOECT environment variable is not set on go1.13 runtime automatically

$ gcloud scheduler jobs create pubsub bqnotify \
  --message-body {} \
  --topic bqnotify \
  --schedule '0 0 * * *'
```

## Notification result example

![bqnotify.jpg](../bqnotify.jpg)
