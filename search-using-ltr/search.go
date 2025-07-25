package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
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

type result struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	ReleaseYear string  `json:"release_year"`
	Score       float64 `json:"_score"`
}

func (c *Client) Search(ctx context.Context, index string, queryBuilder QueryBuilder) ([]result, error) {
	query, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := c.typedClient.Search().
		Index(index).
		Raw(strings.NewReader(string(query))).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}

	var results []result
	for _, hit := range res.Hits.Hits {
		var r result
		if err := json.Unmarshal(hit.Source_, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal search result: %w", err)
		}

		r.Score = float64(*hit.Score_)
		results = append(results, r)
	}
	return results, nil
}
