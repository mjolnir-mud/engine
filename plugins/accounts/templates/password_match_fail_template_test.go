package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordMatchFailTemplate_Name(t *testing.T) {
	assert.Equal(t, PasswordMatchFailTemplate.Name(), "password_match_fail")
}

func TestPasswordMatchFailTemplate_Style(t *testing.T) {
	assert.Equal(t, PasswordMatchFailTemplate.Style(), "default")
}

func TestPasswordMatchFailTemplate_Render(t *testing.T) {
	s, err := PasswordMatchFailTemplate.Render(nil)

	assert.NoError(t, err)
	assert.Equal(t, s, "The password and the confirmation do not match.")
}
