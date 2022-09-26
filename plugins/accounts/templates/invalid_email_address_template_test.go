package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidEmailAddressTemplate_Name(t *testing.T) {
	assert.Equal(t, InvalidEmailAddressTemplate.Name(), "invalid_email_address")
}

func TestInvalidEmailAddressTemplate_Style(t *testing.T) {
	assert.Equal(t, InvalidEmailAddressTemplate.Style(), "default")
}

func TestInvalidEmailAddressTemplate_Render(t *testing.T) {
	s, err := InvalidEmailAddressTemplate.Render("notAnEmail")

	assert.NoError(t, err)
	assert.Equal(t, s, "The provided email address 'notAnEmail' is invalid.")
}
