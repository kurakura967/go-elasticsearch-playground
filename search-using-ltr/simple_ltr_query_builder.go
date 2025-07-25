package main

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// SimpleLTRQueryBuilder builds simple LTR queries with minimal configuration
type SimpleLTRQueryBuilder struct {
	BaseQuery   *types.Query
	RescorePart *Rescore
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

// Build implements the QueryBuilder interface
func (r SimpleLTRQueryBuilder) Build() ([]byte, error) {
	req := &search.Request{
		Query: r.BaseQuery,
	}
	base, err := json.Marshal(req)
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

func NewSimpleLTRQueryBuilder(baseQuery *types.Query, model_name string) SimpleLTRQueryBuilder {
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
	return SimpleLTRQueryBuilder{
		BaseQuery:   baseQuery,
		RescorePart: &rescore,
	}
}

func (s *SimpleLTRQueryBuilder) WithWindowSize(size int) *SimpleLTRQueryBuilder {
	s.RescorePart.WindowSize = size
	return s
}

func (s *SimpleLTRQueryBuilder) WithParams(params map[string]interface{}) *SimpleLTRQueryBuilder {
	sltr := s.RescorePart.Query.RescoreQuery["sltr"]
	if sltr.Params == nil {
		sltr.Params = make(map[string]interface{})
	}
	for k, v := range params {
		sltr.Params[k] = v
	}
	s.RescorePart.Query.RescoreQuery["sltr"] = sltr
	return s
}
