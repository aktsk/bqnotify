package runner

import (
	"bytes"
	"os"

	"github.com/aktsk/bqnotify/lib/bq"
	"github.com/aktsk/bqnotify/lib/config"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/sync/errgroup"
)

// Run coordinates functions of bqnotify
func Run(file string) error {
	conf, err := config.Parse(file)
	if err != nil {
		return err
	}

	eg := errgroup.Group{}
	for _, query := range conf.Queries {
		query := query // capture variable for goroutine
		if query.Slack == nil {
			query.Slack = conf.Slack
		}

		eg.Go(func() error {
			return run(conf.Project, query)
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func run(project string, query config.Query) error {
	headers, rows, err := bq.Query(project, query)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}

	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)

	table.SetHeader(headers)

	for _, v := range rows {
		table.Append(v)
	}

	table.Render()

	query.Slack.URL = os.Getenv("SLACK_WEBHOOK_URL")
	query.Slack.Notify("```\n" + buf.String() + "```")

	return nil
}
