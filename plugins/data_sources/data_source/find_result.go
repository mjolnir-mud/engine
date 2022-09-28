package data_source

import (
	"fmt"
)

type FindResult struct {
	Type     string
	Id       string
	Metadata map[string]interface{}
	Record   map[string]interface{}
}

// NewFindResult creates a new FindResult given the type, id, and result.
func NewFindResult(args ...interface{}) *FindResult {
	t, ok := args[0].(string)

	if !ok {
		panic(fmt.Sprintf("expected string for type, got %T", args[0]))
	}

	i, ok := args[1].(string)

	if !ok {
		panic(fmt.Sprintf("expected string for id, got %T", args[1]))
	}

	r, ok := args[2].(map[string]interface{})

	if !ok {
		panic(fmt.Sprintf("expected map[string]interface{} for result, got %T", args[2]))
	}

	m, ok := args[3].(map[string]interface{})

	if !ok {
		panic(fmt.Sprintf("expected map[string]interface{} for metadata, got %T", args[3]))
	}

	return &FindResult{
		Type:     t,
		Id:       i,
		Metadata: m,
		Record:   r,
	}
}
