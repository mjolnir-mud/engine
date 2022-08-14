package session

import (
	testing2 "github.com/mjolnir-mud/engine/plugins/world/pkg/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStart(t *testing.T) {
	testing2.Setup()
	defer testing2.Teardown()

	err := Start("test")

	assert.NoError(t, err)
}
