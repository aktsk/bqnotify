package bq

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/mizzy/bqnotify/lib/config"
	"google.golang.org/api/iterator"
)

func Query(conf *config.Config) ([]string, [][]string, error) {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, conf.Project)
	if err != nil {
		return nil, nil, err
	}

	q := client.Query(conf.SQL)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, nil, err
	}

	headers := []string{}
	for _, s := range it.Schema {
		headers = append(headers, s.Name)
	}

	rows := [][]string{[]string{}}

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
