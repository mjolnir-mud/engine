package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptNewPasswordTemplate_Name(t *testing.T) {
	assert.Equal(t, PromptNewPasswordTemplate.Name(), "prompt_new_password")
}

func TestPromptNewPasswordTemplate_Style(t *testing.T) {
	assert.Equal(t, PromptNewPasswordTemplate.Style(), "default")
}

func TestPromptNewPasswordTemplate_Render(t *testing.T) {
	s, err := PromptNewPasswordTemplate.Render(nil)

	assert.NoError(t, err)
	assert.Equal(t, s, "Enter a password:")
}
