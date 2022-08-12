package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var beforeStartCalled = make(chan bool)
var afterStartCalled = make(chan bool)

var beforeStopCalled = make(chan bool)
var afterStopCalled = make(chan bool)

func TestStart(t *testing.T) {
	RegisterBeforeStartCallback(func() {
		go func() {
			beforeStartCalled <- true
		}()
	})

	RegisterAfterStartCallback(func() {
		go func() {
			afterStartCalled <- true
		}()
	})

	Start("test")

	<-beforeStartCalled
	<-afterStartCalled

	assert.True(t, IsRunning())
}

func TestStop(t *testing.T) {
	RegisterBeforeStopCallback(func() {
		go func() {
			beforeStopCalled <- true
		}()
	})

	RegisterAfterStopCallback(func() {
		go func() {
			afterStopCalled <- true
		}()
	})

	Start("test")

	Stop()
	<-beforeStopCalled
	<-afterStopCalled
	assert.True(t, true)
}
