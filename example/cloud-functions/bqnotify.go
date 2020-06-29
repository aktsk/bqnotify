package example

import (
	"context"

	"github.com/aktsk/bqnotify/lib/runner"
)

// PubSubMessage is needed to run on Cloud Functions
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// BqNotify is a main function to run bqnotify on Cloud Functions
func BqNotify(ctx context.Context, m PubSubMessage) error {
	err := runner.Run("./serverless_function_source_code/config.yaml")
	if err != nil {
		return err
	}
	return nil
}
