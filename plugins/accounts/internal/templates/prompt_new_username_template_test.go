package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptNewUsernameTemplate_Name(t *testing.T) {
	assert.Equal(t, PromptNewUsernameTemplate.Name(), "prompt_new_username")
}

func TestPromptNewUsernameTemplate_Style(t *testing.T) {
	assert.Equal(t, PromptNewUsernameTemplate.Style(), "default")
}

func TestPromptNewUsernameTemplate_Render(t *testing.T) {
	s, err := PromptNewUsernameTemplate.Render(nil)

	assert.NoError(t, err)
	assert.Equal(t, s, "Enter a username:")
}
