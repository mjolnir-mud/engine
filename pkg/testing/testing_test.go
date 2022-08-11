package testing

import (
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/stretchr/testify/assert"
)

func TestSetupAndTeardown(t *testing.T) {
	ch := Setup()
	<-ch

	err := engine.Ping()

	assert.Nil(t, err)
	Teardown()

	err = engine.Ping()
	assert.NotNil(t, err)
}
