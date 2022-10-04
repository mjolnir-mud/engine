package templates

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/template_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/theme_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/themes/default_theme"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	engineTesting.RegisterSetupCallback("templates", func() {
		template_registry.Start()
		theme_registry.Start()
		engine.RegisterAfterServiceStartCallback("world", func() {
			theme_registry.Register(default_theme.Theme)
			Register()
		})
	})

	engineTesting.Setup("world")
}

func teardown() {
	theme_registry.Stop()
	template_registry.Stop()
	engineTesting.Teardown()
}

func TestOrderedListTemplate_Name(t *testing.T) {
	assert.Equal(t, "ordered_list", OrderedList.Name())
}

func TestOrderedListTemplate_Style(t *testing.T) {
	assert.Equal(t, "default", OrderedList.Style())
}

func TestOrderedListTemplate_Render(t *testing.T) {
	setup()
	defer teardown()

	s, err := OrderedList.Render([]string{
		"one",
		"two",
		"three",
	})

	assert.Equal(t, nil, err)

	assert.Equal(t, "    \x1b[1m1. \x1b[0mone\n    \x1b[1m2. \x1b[0mtwo\n    \x1b[1m3. \x1b[0mthree", s)
}
