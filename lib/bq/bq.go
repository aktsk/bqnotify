package bq

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/aktsk/bqnotify/lib/config"
	"google.golang.org/api/iterator"
)

type Result struct {
	Headers   []string
	Rows      [][]string
	DatasetID string
	TableID   string
}

// Query runs queries to BigQuery and return results
func Query(project string, query config.Query) (*Result, error) {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	q := client.Query(query.SQL)

	var datasetID string
	var tableID string

	// Write the query result to the table specified in config.yaml
	if query.ResultTable != nil {
		datasetID = query.ResultTable.DatasetID
		tableID = fmt.Sprintf("%s%s", query.ResultTable.TableIDPrefix, time.Now().Format("20060102150405"))

		// Create dataset if it does not exist
		metadata, _ := client.Dataset(datasetID).Metadata(ctx)
		if metadata == nil {
			err := client.Dataset(datasetID).Create(ctx, &bigquery.DatasetMetadata{})
			if err != nil {
				return nil, err
			}
		}

		// Set the destination table to write the query result
		q.QueryConfig.Dst = client.Dataset(datasetID).Table(tableID)
	}

	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}

	it, err := job.Read(ctx)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		values := []string{}
		for _, v := range bqValues {
			values = append(values, fmt.Sprint(v))
		}
		rows = append(rows, values)
	}

	// Set expiration time to the result table
	if query.ResultTable != nil {
		tableRef := client.Dataset(datasetID).Table(tableID)

		meta, err := tableRef.Metadata(ctx)
		if err != nil {
			return nil, err
		}

		if query.ResultTable.ExpirationInDays == 0 {
			query.ResultTable.ExpirationInDays = 30
		}

		update := bigquery.TableMetadataToUpdate{
			ExpirationTime: time.Now().Add(time.Duration(query.ResultTable.ExpirationInDays*24) * time.Hour),
		}

		_, err = tableRef.Update(ctx, update, meta.ETag)
		if err != nil {
			return nil, err
		}
	}

	return &Result{
		Headers:   headers,
		Rows:      rows,
		DatasetID: datasetID,
		TableID:   tableID,
	}, nil
}
