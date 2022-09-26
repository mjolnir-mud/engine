package account

import (
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	engineTesting.Setup("world", func() {
		ecsTesting.Setup()
	})
}

func teardown() {
	ecsTesting.Teardown()
	engineTesting.Teardown()
}

func TestSystem_Name(t *testing.T) {
	setup()
	defer teardown()
	assert.Equal(t, "account", System.Name())
}

func TestSystem_Component(t *testing.T) {
	setup()
	defer teardown()
	assert.Equal(t, "_", System.Component())
}

func TestValidateEmail(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ValidateEmail("test@email.com"))

	assert.Error(t, ValidateEmail("test"))
}

func TestValidatePassword(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ValidatePassword("test", "test@test.com", "validPassword"))
	assert.Error(t, ValidatePassword("test", "test@test.com", "invalid"))
	assert.Error(t, ValidatePassword("validPassword", "test@test.com", "validPassword"))
	assert.Error(t, ValidatePassword("test", "test@test.com", "test@test.com"))
}

func TestValidateUsername(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ValidateUsername("test"))
	assert.Error(t, ValidateUsername("invalid@"))
}
