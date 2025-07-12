package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Client struct {
	baseClient  *elasticsearch.Client
	typedClient *elasticsearch.TypedClient
}

func NewClient(cfg elasticsearch.Config) (*Client, error) {
	typedClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}

	baseClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseClient:  baseClient,
		typedClient: typedClient,
	}, nil
}

func (c *Client) CreateIndex(ctx context.Context, index string) error {
	exist, err := c.typedClient.Indices.Exists(index).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if index exists: %w", err)
	}
	if exist {
		_, err := c.typedClient.Indices.Delete(index).Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete existing index: %w", err)
		}
	}

	const body = `{
		"settings": {
			"refresh_interval": "60s",
			"number_of_shards": 1,
			"auto_expand_replicas": "0-all"
		},
		"mappings": {
			"dynamic": "strict",
			"properties": {
				"title": {
					"type": "text"
				},
				"author": {
					"type": "text"
				}
			}
		}
	}`
	_, err = c.typedClient.Indices.Create(index).Raw(strings.NewReader(body)).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	return nil
}

func (c *Client) SingleInsert(ctx context.Context, index string, docs []map[string]interface{}) error {
	for i, doc := range docs {
		_, err := c.typedClient.Index(index).
			Id(fmt.Sprintf("%d", i+1)).
			Request(doc).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to insert document %d: %w", i+1, err)
		}
	}
	return nil
}

func (c *Client) BulkInsert(ctx context.Context, index string, docs []map[string]interface{}) error {
	bulkCfg := esutil.BulkIndexerConfig{
		Client: c.baseClient,
		Index:  index,
	}

	indexer, err := esutil.NewBulkIndexer(bulkCfg)
	if err != nil {
		return fmt.Errorf("failed to create bulk indexer: %w", err)
	}
	for i, doc := range docs {
		data, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document %d: %w", i+1, err)
		}
		err = indexer.Add(
			ctx,
			esutil.BulkIndexerItem{
				Index:      index,
				Action:     "index",
				DocumentID: fmt.Sprintf("%d", i+1),
				Body:       strings.NewReader(string(data)),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to add document %d to bulk indexer: %w", i+1, err)
		}
	}
	if err := indexer.Close(ctx); err != nil {
		return fmt.Errorf("failed to close bulk indexer: %w", err)
	}
	return nil
}
