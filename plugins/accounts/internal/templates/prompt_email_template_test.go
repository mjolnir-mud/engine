package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptEmailTemplate_Name(t *testing.T) {
	assert.Equal(t, PromptEmailTemplate.Name(), "prompt_email")
}

func TestPromptEmailTemplate_Style(t *testing.T) {
	assert.Equal(t, PromptEmailTemplate.Style(), "default")
}

func TestPromptEmailTemplate_Render(t *testing.T) {
	s, err := PromptEmailTemplate.Render(nil)

	assert.NoError(t, err)
	assert.Equal(t, s, "Enter your email:")
}
