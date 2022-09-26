package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptNewEmailTemplate_Name(t *testing.T) {
	assert.Equal(t, PromptNewEmailTemplate.Name(), "prompt_new_email")
}

func TestPromptNewEmailTemplate_Style(t *testing.T) {
	assert.Equal(t, PromptNewEmailTemplate.Style(), "default")
}

func TestPromptNewEmailTemplate_Render(t *testing.T) {
	s, err := PromptNewEmailTemplate.Render(nil)

	assert.NoError(t, err)
	assert.Equal(t, s, "Enter an email:")
}
