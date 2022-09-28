package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetadataRequiredError_Error(t *testing.T) {
	assert.Equal(t, "entity with id 123 did not return an entity with the type set in the metadata", MetadataRequiredError{ID: "123"}.Error())
}
