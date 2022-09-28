package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntityNotFoundError_Error(t *testing.T) {
	assert.Equal(t, "entity 123 not found", EntityNotFoundError{ID: "123"}.Error())
}
