package example

import (
	"context"

	"github.com/mizzy/bqnotify/lib/runner"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func BqNotify(ctx context.Context, m PubSubMessage) error {
	err := runner.Run()
	if err != nil {
		return err
	}
	return nil
}
