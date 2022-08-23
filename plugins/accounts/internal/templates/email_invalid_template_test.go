package templates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailInvalidTemplate_Name(t *testing.T) {
	assert.Equal(t, "email_invalid", EmailInvalidTemplate.Name())
}

func TestEmailInvalidTemplate_Style(t *testing.T) {
	assert.Equal(t, "default", EmailInvalidTemplate.Style())
}

func TestEmailInvalidTemplate_Render(t *testing.T) {
	s, err := EmailInvalidTemplate.Render("notan email")

	assert.Equal(t, nil, err)

	assert.Equal(t, "'notan email' is not a valid email address.", s)
}
