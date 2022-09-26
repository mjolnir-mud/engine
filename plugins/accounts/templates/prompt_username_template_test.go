package templates

import (
	"github.com/mjolnir-mud/engine/plugins/accounts/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptUsernameTemplate_Render(t *testing.T) {
	test.Setup()
	defer test.Teardown()

	v, err := PromptUsernameTemplate.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "Enter your username, or type '\u001B[1mcreate\u001B[0m' to create a new account:", v)
}
