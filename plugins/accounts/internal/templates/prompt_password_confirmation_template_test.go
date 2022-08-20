package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptPasswordConfirmationTemplate_Name(t *testing.T) {
	assert.Equal(t, PromptPasswordConfirmationTemplate.Name(), "prompt_password_confirmation")
}

func TestPromptPasswordConfirmationTemplate_Style(t *testing.T) {
	assert.Equal(t, PromptPasswordConfirmationTemplate.Style(), "default")
}

func TestPromptPasswordConfirmationTemplate_Render(t *testing.T) {
	s, err := PromptPasswordConfirmationTemplate.Render(nil)

	assert.NoError(t, err)
	assert.Equal(t, s, "Confirm your password:")
}
