package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {

	beforeStartCalled := make(chan bool)
	afterStartCalled := make(chan bool)

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
	defer Stop()

	<-beforeStartCalled
	<-afterStartCalled

	assert.True(t, IsRunning())
}

func TestStop(t *testing.T) {
	beforeStopCalled := make(chan bool)
	afterStopCalled := make(chan bool)

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
