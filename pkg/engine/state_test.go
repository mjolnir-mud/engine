package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	Init("test", []Plugin{&testPlugin{}})

	assert.Equal(t, "test", Name())

}
