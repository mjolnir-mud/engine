package templates

import (
	"github.com/mjolnir-mud/engine/plugins/templates/internal/plugin"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/template_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/theme_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/template"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/theme"
)

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

var Plugin = plugin.Plugin
