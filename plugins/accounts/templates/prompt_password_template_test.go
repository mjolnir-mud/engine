package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptPasswordTemplate_Name(t *testing.T) {
	assert.Equal(t, "prompt_password", PromptPasswordTemplate.Name())
}

func TestPromptPasswordTemplate_Style(t *testing.T) {
	assert.Equal(t, "default", PromptPasswordTemplate.Style())
}

func TestPromptPasswordTemplate_Render(t *testing.T) {
	v, err := PromptPasswordTemplate.Render(nil)
	assert.NoError(t, err)
	assert.Equal(t, "Enter your password:", v)
}
