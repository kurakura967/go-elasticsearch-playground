package main

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// QueryType represents the type of query builder to use
type QueryType string

const (
	QueryTypeLTR       QueryType = "ltr"
	QueryTypeSimpleLTR QueryType = "simple-ltr"
	QueryTypeStringLTR QueryType = "string-ltr"
)

// CreateLTRQueryBuilder is a factory method that creates the appropriate QueryBuilder
func CreateLTRQueryBuilder(queryType QueryType, baseQuery *types.Query, keyword, model string) QueryBuilder {
	switch queryType {
	case QueryTypeLTR:
		return NewLTRQueryBuilder(baseQuery, model)
	case QueryTypeSimpleLTR:
		return NewSimpleLTRQueryBuilder(baseQuery, model)
	case QueryTypeStringLTR:
		return NewStringLTRQueryBuilder(keyword, model)
	default:
		// Default to LTRQueryBuilder
		return NewLTRQueryBuilder(baseQuery, model)
	}
}

func main() {
	client, err := NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		fmt.Printf("Error creating Elasticsearch client: %v\n", err)
		return
	}

	ctx := context.Background()

	// Create base query
	baseQuery := &types.Query{
		Bool: &types.BoolQuery{
			Must: []types.Query{
				{MatchAll: types.NewMatchAllQuery()},
			},
			Filter: []types.Query{
				{
					Match: map[string]types.MatchQuery{
						"title": {Query: "batman"},
					},
				},
			},
		},
	}

	// Example 1: Using SimpleLTRQueryBuilder via factory
	fmt.Println("=== Using SimpleLTRQueryBuilder ===")
	simpleBuilder := CreateLTRQueryBuilder(QueryTypeSimpleLTR, baseQuery, "batman", "latest")
	if sb, ok := simpleBuilder.(*SimpleLTRQueryBuilder); ok {
		sb.WithWindowSize(500).
			WithParams(map[string]interface{}{
				"hoge": "fuga",
			})
	}
	res, err := client.Search(ctx, "tmdb", simpleBuilder)
	if err != nil {
		fmt.Printf("Error executing search: %v\n", err)
		return
	}
	for i, r := range res {
		if i >= 3 {
			break
		}
		fmt.Printf("Found result: ID=%s, Title=%s, Release Year=%s, Score=%.2f\n",
			r.Id, r.Title, r.ReleaseYear, r.Score)
	}

	fmt.Println()

	// Example 2: Using LTRQueryBuilder via factory
	fmt.Println("=== Using LTRQueryBuilder ===")
	ltrBuilder := CreateLTRQueryBuilder(QueryTypeLTR, baseQuery, "batman", "latest")

	// Optionally configure the LTR builder
	if lb, ok := ltrBuilder.(*LTRQueryBuilder); ok {
		lb.WithWindowSize(500)
	}

	res, err = client.Search(ctx, "tmdb", ltrBuilder)
	if err != nil {
		fmt.Printf("Error executing search: %v\n", err)
		return
	}
	for i, r := range res {
		if i >= 3 {
			break
		}
		fmt.Printf("Found result: ID=%s, Title=%s, Release Year=%s, Score=%.2f\n",
			r.Id, r.Title, r.ReleaseYear, r.Score)
	}

	fmt.Println()

	// Example 3: Using StringLTRQueryBuilder via factory
	fmt.Println("=== Using StringLTRQueryBuilder ===")
	stringBuilder := CreateLTRQueryBuilder(QueryTypeStringLTR, nil, "batman", "latest")

	res, err = client.Search(ctx, "tmdb", stringBuilder)
	if err != nil {
		fmt.Printf("Error executing search: %v\n", err)
		return
	}
	for i, r := range res {
		if i >= 3 {
			break
		}
		fmt.Printf("Found result: ID=%s, Title=%s, Release Year=%s, Score=%.2f\n",
			r.Id, r.Title, r.ReleaseYear, r.Score)
	}
}
