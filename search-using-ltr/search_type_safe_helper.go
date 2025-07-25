package main

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// ヘルパー関数: よく使われるクエリパターンを簡単に構築

// CreateMatchQuery は単純なマッチクエリを作成
func CreateMatchQuery(field, value string) *types.Query {
	return &types.Query{
		Match: map[string]types.MatchQuery{
			field: {Query: value},
		},
	}
}

// CreateBoolWithFilterQuery はフィルター付きのBoolクエリを作成
func CreateBoolWithFilterQuery(filters ...types.Query) *types.Query {
	return &types.Query{
		Bool: &types.BoolQuery{
			Must:   []types.Query{{MatchAll: types.NewMatchAllQuery()}},
			Filter: filters,
		},
	}
}

// CreateMultiMatchQuery は複数フィールドに対するマッチクエリを作成
func CreateMultiMatchQuery(query string, fields ...string) *types.Query {
	return &types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query:  query,
			Fields: fields,
		},
	}
}

// CreateComplexBoolQuery はより複雑なBoolクエリを作成
func CreateComplexBoolQuery(must, should, filter []types.Query, minimumShouldMatch *string) *types.Query {
	return &types.Query{
		Bool: &types.BoolQuery{
			Must:               must,
			Should:             should,
			Filter:             filter,
			MinimumShouldMatch: minimumShouldMatch,
		},
	}
}