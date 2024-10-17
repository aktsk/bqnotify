package bq

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/aktsk/bqnotify/lib/config"
	"google.golang.org/api/iterator"
)

// Query runs queries to BigQuery and return results
func Query(project string, query config.Query) ([]string, [][]string, error) {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, nil, err
	}

	defer client.Close()

	q := client.Query(query.SQL)

	job, err := q.Run(ctx)
	if err != nil {
		return nil, nil, err
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return nil, nil, err
	}
	if err := status.Err(); err != nil {
		return nil, nil, err
	}

	it, err := job.Read(ctx)
	if err != nil {
		return nil, nil, err
	}

	headers := []string{}
	for _, s := range it.Schema {
		headers = append(headers, s.Name)
	}

	var rows [][]string

	for {
		var bqValues []bigquery.Value

		err := it.Next(&bqValues)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, nil, err
		}

		values := []string{}
		for _, v := range bqValues {
			values = append(values, fmt.Sprint(v))
		}
		rows = append(rows, values)
	}

	return headers, rows, nil
}
