package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIDRequiredError_Error(t *testing.T) {
	assert.Equal(t, "entity ID required", IDRequiredError{}.Error())
}
