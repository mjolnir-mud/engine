package template_registry

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/theme_registry"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testTemplate struct{}

func (t *testTemplate) Render(data interface{}) (string, error) {
	return "", nil
}

func (t *testTemplate) Name() string {
	return "test"
}

func (t *testTemplate) Style() string {
	return "test"
}

type testTheme struct{}

func (t *testTheme) Name() string {
	return "default"
}

func (t *testTheme) GetStyleFor(_ string) lipgloss.Style {
	return lipgloss.Style{}
}

func (t *testTheme) DefaultStyle() lipgloss.Style {
	return lipgloss.Style{}
}

func setup() {
	theme_registry.Start()
	theme_registry.Register(&testTheme{})
	Start()
}

func teardown() {
	Stop()
}

func TestRegister(t *testing.T) {
	setup()
	defer teardown()

	Register(&testTemplate{})

	assert.Len(t, templates, 1)
}

func TestRender(t *testing.T) {
	Start()
	defer Stop()

	Register(&testTemplate{})

	_, err := Render("test", nil)

	assert.NoError(t, err)
}
