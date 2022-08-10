package testing

import (
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/stretchr/testify/assert"
)

func TestSetupAndTeardown(t *testing.T) {
	ch := make(chan bool)

	go Setup(func() {
		ch <- true
	})

	<-ch

	err := engine.Ping().Err()

	assert.Nil(t, err)
	Teardown()

	err = engine.Ping().Err()
	assert.NotNil(t, err)
}
