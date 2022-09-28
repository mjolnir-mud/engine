package data_source

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindResults_All(t *testing.T) {
	fr := NewFindResults([]*FindResult{
		{Type: "test", Id: "1", Record: map[string]interface{}{}, Metadata: map[string]interface{}{}},
	})

	assert.Equal(t, 1, len(fr.All()))
}

func TestFindResults_Get(t *testing.T) {
	fr := NewFindResults([]*FindResult{
		{Type: "test", Id: "1", Record: map[string]interface{}{}, Metadata: map[string]interface{}{}},
	})

	assert.Equal(t, "1", fr.Get("1").Id)
}

func TestFindResults_Len(t *testing.T) {
	fr := NewFindResults([]*FindResult{
		{Type: "test", Id: "1", Record: map[string]interface{}{}, Metadata: map[string]interface{}{}},
	})

	assert.Equal(t, 1, fr.Len())
}
