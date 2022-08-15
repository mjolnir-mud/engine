package login_controller

import (
	"github.com/magiconair/properties/assert"
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/world"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	testing2 "testing"
)

func setup() {
	engine.RegisterPlugin(world.Plugin)
	testing.Setup()
}

func teardown() {
	testing.Teardown()
}

func TestController_Name(t *testing2.T) {
	setup()
	defer teardown()

	c := controller{}

	assert.Equal(t, "login", c.Name())
}

func TestController_Start(t *testing2.T) {
	setup()
	defer teardown()

	receivedLine := make(chan string)

	sub := engine.Subscribe(events.SendLineEvent{}, "sess", func(e interface{}) {
		go func() { receivedLine <- e.(events.SendLineEvent).Line }()
	})

	defer sub.Stop()

	c := controller{}

	err := c.Start("sess")

	assert.Equal(t, nil, err)

	line := <-receivedLine

	assert.Equal(t, "Username: ", line)
}
