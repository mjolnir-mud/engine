package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidDataSourceError_Error(t *testing.T) {
	assert.Equal(t, "data source 123 does not exist", InvalidDataSourceError{Source: "123"}.Error())
}
