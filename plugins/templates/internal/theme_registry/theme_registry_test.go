package theme_registry

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testTheme struct{}

func (t testTheme) Name() string {
	return "testing"
}

func (t testTheme) DefaultStyle() lipgloss.Style {
	return lipgloss.Style{}
}

func (t testTheme) GetStyleFor(_ string) lipgloss.Style {
	return lipgloss.Style{}
}

func setup() {
	Start()
}

func teardown() {
	Stop()
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()

	Register(testTheme{})

	assert.Len(t, themes, 1)
}

func TestGetTheme(t *testing.T) {
	setup()
	defer teardown()

	Register(testTheme{})

	thm, err := GetTheme("testing")

	assert.NoError(t, err)
	assert.Equal(t, testTheme{}.Name(), thm.Name())

	_, err = GetTheme("test2")

	assert.Error(t, err)

}
