package data_source

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFindResult(t *testing.T) {
	fr := NewFindResult("fake", "fake", map[string]interface{}{}, map[string]interface{}{})

	assert.NotNil(t, fr)
}
