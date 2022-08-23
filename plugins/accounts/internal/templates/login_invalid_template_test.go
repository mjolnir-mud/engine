package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoginInvalidTemplate_Name(t *testing.T) {
	assert.Equal(t, "login_invalid", LoginInvalidTemplate.Name())
}

func TestLoginInvalidTemplate_Style(t *testing.T) {
	assert.Equal(t, "default", LoginInvalidTemplate.Style())
}

func TestLoginInvalidTemplate_Render(t *testing.T) {
	str, err := LoginInvalidTemplate.Render(nil)

	assert.Nil(t, err)
	assert.Equal(
		t,
		"An account with that username and password combination was not found.",
		str,
	)
}
