package runner

import (
	"bytes"
	"os"

	"github.com/aktsk/bqnotify/lib/bq"
	"github.com/aktsk/bqnotify/lib/config"
	"github.com/olekukonko/tablewriter"
)

func Run() error {
	conf, err := config.Parse()
	if err != nil {
		return err
	}

	headers, values, err := bq.Query(conf)
	if err != nil {
		return err
	}

	if len(values) == 0 {
		return nil
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)

	table.SetHeader(headers)

	for _, v := range values {
		table.Append(v)
	}

	table.Render()

	conf.Slack.URL = os.Getenv("SLACK_WEBHOOK_URL")
	conf.Slack.Notify("```\n" + buf.String() + "```")

	return nil
}
