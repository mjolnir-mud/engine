package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsernameTakenTemplate_Name(t *testing.T) {
	assert.Equal(t, UsernameTakenTemplate.Name(), "username_taken")
}

func TestUsernameTakenTemplate_Style(t *testing.T) {
	assert.Equal(t, UsernameTakenTemplate.Style(), "default")
}

func TestUsernameTakenTemplate_Render(t *testing.T) {
	s, err := UsernameTakenTemplate.Render("testing")

	assert.NoError(t, err)
	assert.Equal(t, s, "The username 'testing' is already taken.")
}
