package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	Start()
}

func teardown() {
	Stop()
}

type testController struct{}

func (c testController) Name() string {
	return "test"
}

func (c testController) Start(_ string) error {
	return nil
}

func (c testController) Resume(_ string) error {
	return nil
}

func (c testController) Stop(_ string) error {
	return nil
}

func (c testController) HandleInput(_ string, _ string) error {
	return nil
}

func TestStart(t *testing.T) {
	setup()
	defer teardown()

	Start()

	assert.NotNil(t, controllers)
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()

	Register(testController{})

	assert.Len(t, controllers, 1)
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	Register(testController{})

	c, err := Get("test")

	assert.Nil(t, err)
	assert.NotNil(t, c)
}
