# Example

## How to run with Cloud Functions/Cloud Scheduler

```sh
$ gcloud pubsub topics create bqnotify

$ gcloud functions deploy BqNotify \
    --region asia-northeast1 \
    --runtime go113 \
    --trigger-topic bqnotify \
    --set-env-vars SLACK_WEBHOOK_URL=$SLACK_WEBHOOK_URL

$ gcloud scheduler jobs create pubsub bqnotify \
  --message-body {} \
  --topic bqnotify \
  --schedule '0 0 * * *'
```
