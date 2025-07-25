package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
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

func buildJsonQuery(keyword, model_name string) string {
	return `{
  "query": {
    "bool": {
      "must": {"match_all": {}},
      "filter": {"match": {"title": "` + keyword + `"}}
    }
  },
  "rescore": {
    "window_size": 1000,
    "query": {
      "rescore_query": {
        "sltr": {
          "params": {},
          "model": "` + model_name + `"
        }
      }
    }
  }
}`
}

type result struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	ReleaseYear string  `json:"release_year"`
	Score       float64 `json:"_score"`
}

func (c *Client) Search(ctx context.Context, index string, query string) ([]result, error) {
	res, err := c.typedClient.Search().
		Index(index).
		Raw(strings.NewReader(query)).
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

type SLTR struct {
	Params map[string]interface{} `json:"params"`
	Model  string                 `json:"model"`
}

type RescoreQuery struct {
	RescoreQuery map[string]SLTR `json:"rescore_query"`
}

type Rescore struct {
	WindowSize int          `json:"window_size"`
	Query      RescoreQuery `json:"query"`
}

type CustomSearchRequest struct {
	*search.Request
	RescorePart *Rescore
}

func (r CustomSearchRequest) MarshalJSON() ([]byte, error) {
	base, err := json.Marshal(r.Request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal base request: %w", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base request: %w", err)
	}

	if r.RescorePart != nil {
		m["rescore"] = r.RescorePart
	}

	return json.Marshal(m)
}

func NewCustomSearchRequest(keyword, model_name string) CustomSearchRequest {

	q := types.Query{
		Bool: &types.BoolQuery{
			Must: []types.Query{
				{
					MatchAll: types.NewMatchAllQuery(),
				},
			},
			Filter: []types.Query{
				{
					Match: map[string]types.MatchQuery{
						"title": {
							Query: keyword,
						},
					},
				},
			},
		},
	}
	rescore := Rescore{
		WindowSize: 1000,
		Query: RescoreQuery{
			RescoreQuery: map[string]SLTR{
				"sltr": {
					Params: map[string]interface{}{},
					Model:  model_name,
				},
			},
		},
	}
	r := search.Request{
		Query: &q,
	}
	return CustomSearchRequest{
		Request:     &r,
		RescorePart: &rescore,
	}
}
