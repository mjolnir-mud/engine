package templates

import (
	"github.com/mjolnir-mud/engine/plugins/templates/internal/template_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/theme_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/default_theme"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/template"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/theme"
)

type templatePlugin struct {
	themes    map[string]theme.Theme
	templates map[string]template.Template
}

func (p templatePlugin) Name() string {
	return "templates"
}

func (p templatePlugin) Registered() error {
	theme_registry.Start()
	template_registry.Start()

	theme_registry.Register(default_theme.Theme)
	return nil
}

// RegisterTheme registers a theme with the theme registry.
func RegisterTheme(t theme.Theme) {
	theme_registry.Register(t)
}

// RegisterTemplate registers a template with the template registry.
func RegisterTemplate(t template.Template) {
	template_registry.Register(t)
}

// GetTheme returns a theme with the given name. If the theme is not found, an error is returned.
func GetTheme(name string) (theme.Theme, error) {
	return theme_registry.GetTheme(name)
}

// RenderTemplate renders a template with the given name passing the given data to the template. If the template is not
// found, an error is returned.
func RenderTemplate(name string, ctx interface{}) (string, error) {
	return template_registry.Render(name, ctx)
}

var Plugin = &templatePlugin{}
