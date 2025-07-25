package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	c := NewCustomSearchRequest("batman", "latest")
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(c); err != nil {
		fmt.Printf("Error encoding search request: %v\n", err)
		return
	}

	client, err := NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		fmt.Printf("Error creating Elasticsearch client: %v\n", err)
		return
	}
	ctx := context.Background()
	res, err := client.Search(ctx, "tmdb", buf.String())
	if err != nil {
		fmt.Printf("Error executing search: %v\n", err)
		return
	}
	for _, r := range res {
		fmt.Printf("Found result: ID=%s, Title=%s, Release Year=%s, Score=%.2f\n",
			r.Id, r.Title, r.ReleaseYear, r.Score)

	}
}
