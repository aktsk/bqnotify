package runner

import (
	"bytes"
	"encoding/csv"
	"fmt"
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
	result, err := bq.Query(project, query)
	if err != nil {
		return err
	}

	rows := result.Rows
	if len(rows) == 0 {
		return nil
	}

	headers := result.Headers

	// Human-readable table
	query.Slack.URL = os.Getenv("SLACK_WEBHOOK_URL")
	if query.Slack.URL != "" {
		var buf bytes.Buffer
		table := tablewriter.NewWriter(&buf)

		table.SetHeader(headers)

		for _, v := range rows {
			table.Append(v)
		}

		table.Render()

		err = query.Slack.Notify("```\n" + buf.String() + "```")
		if err != nil {
			return err
		}
	}

	// CSV upload
	if query.Slack.UploadChannelID != "" {
		var csvBuffer bytes.Buffer
		writer := csv.NewWriter(&csvBuffer)

		err = writer.Write(headers)
		if err != nil {
			return err
		}

		for _, v := range rows {
			if err := writer.Write(v); err != nil {
				return err
			}
		}

		writer.Flush()

		if err := writer.Error(); err != nil {
			return nil
		}

		err = query.Slack.Upload(&csvBuffer)
		if err != nil {
			return err
		}
	}

	// Notify the Dataset ID and Table ID of the result table
	if query.Slack.URL != "" && query.ResultTable != nil {
		url := fmt.Sprintf("https://console.cloud.google.com/bigquery?p=%s&d=%s&t=%s&page=table", project, result.DatasetID, result.TableID)
		message := fmt.Sprintf("This query result has been written to the following table:\n\nProject ID: %s\nDataset ID: %s\nTable ID: %s\n\n%s", project, result.DatasetID, result.TableID, url)
		err = query.Slack.Notify(message)
		if err != nil {
			return err
		}
	}

	return nil
}
