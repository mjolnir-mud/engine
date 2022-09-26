package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailTakenTemplate_Name(t *testing.T) {
	assert.Equal(t, EmailTakenTemplate.Name(), "email_taken")
}

func TestEmailTakenTemplate_Style(t *testing.T) {
	assert.Equal(t, EmailTakenTemplate.Style(), "default")
}

func TestEmailTakenTemplate_Render(t *testing.T) {
	s, err := EmailTakenTemplate.Render(nil)

	assert.NoError(t, err)
	assert.Equal(t, s, "That email already belongs to an account.")
}
