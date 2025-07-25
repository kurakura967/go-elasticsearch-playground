package main

// StringLTRQueryBuilder builds LTR queries using string templates
type StringLTRQueryBuilder struct {
	keyword string
	model   string
}

// NewStringLTRQueryBuilder creates a new StringLTRQueryBuilder
func NewStringLTRQueryBuilder(keyword, model string) *StringLTRQueryBuilder {
	return &StringLTRQueryBuilder{
		keyword: keyword,
		model:   model,
	}
}

// Build implements the QueryBuilder interface
func (s *StringLTRQueryBuilder) Build() ([]byte, error) {
	query := `{
  "query": {
    "bool": {
      "must": {"match_all": {}},
      "filter": {"match": {"title": "` + s.keyword + `"}}
    }
  },
  "rescore": {
    "window_size": 1000,
    "query": {
      "rescore_query": {
        "sltr": {
          "params": {},
          "model": "` + s.model + `"
        }
      }
    }
  }
}`
	return []byte(query), nil
}