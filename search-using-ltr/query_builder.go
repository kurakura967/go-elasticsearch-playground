package main

// QueryBuilder interface for building search queries
type QueryBuilder interface {
	Build() ([]byte, error)
}