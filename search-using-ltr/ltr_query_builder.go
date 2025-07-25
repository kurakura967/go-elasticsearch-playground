package main

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// LTRQueryBuilder builds LTR queries with advanced configuration options
type LTRQueryBuilder struct {
	baseQuery   *types.Query
	rescoreConf RescoreConfig
}

type RescoreConfig struct {
	WindowSize int
	Model      string
	Params     map[string]interface{}
}

func NewLTRQueryBuilder(baseQuery *types.Query, model string) *LTRQueryBuilder {
	return &LTRQueryBuilder{
		baseQuery: baseQuery,
		rescoreConf: RescoreConfig{
			WindowSize: 1000,
			Model:      model,
			Params:     map[string]interface{}{},
		},
	}
}

func (b *LTRQueryBuilder) WithWindowSize(size int) *LTRQueryBuilder {
	b.rescoreConf.WindowSize = size
	return b
}

func (b *LTRQueryBuilder) WithParams(params map[string]interface{}) *LTRQueryBuilder {
	for k, v := range params {
		b.rescoreConf.Params[k] = v
	}
	return b
}

func (b *LTRQueryBuilder) Build() ([]byte, error) {
	baseReq := &search.Request{Query: b.baseQuery}

	baseJSON, err := json.Marshal(baseReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal base query: %w", err)
	}

	var queryMap map[string]interface{}
	if err := json.Unmarshal(baseJSON, &queryMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base query: %w", err)
	}

	queryMap["rescore"] = map[string]interface{}{
		"window_size": b.rescoreConf.WindowSize,
		"query": map[string]interface{}{
			"rescore_query": map[string]interface{}{
				"sltr": map[string]interface{}{
					"params": b.rescoreConf.Params,
					"model":  b.rescoreConf.Model,
				},
			},
		},
	}

	return json.Marshal(queryMap)
}
