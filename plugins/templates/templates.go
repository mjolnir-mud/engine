package templates

import (
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/template"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/theme"
	"github.com/rs/zerolog/log"
)

type templatePlugin struct {
	themes    map[string]theme.Theme
	templates map[string]template.Template
}

func (p templatePlugin) Name() string {
	return "templates"
}

func (p templatePlugin) Start() error {
	return nil
}

func RegisterTheme(t theme.Theme) {
	Plugin.themes[t.Name()] = t
}

func GetTheme(name string) theme.Theme {
	return Plugin.themes[name]
}

func RegisterTemplate(t template.Template) {
	Plugin.templates[t.Name()] = t
}

func RenderTemplate(name string, ctx interface{}) string {
	// get the template
	t := Plugin.templates[name]

	// if the template doesn't exist, return an error
	if t == nil {
		log.Error().Str("template", name).Msg("template not found")
		return "Something went terribly wrong."
	}

	// get the theme and style
	thm := GetTheme("default")
	style := thm.GetStyleFor(t.Style())

	text, err := t.Render(ctx)

	if err != nil {
		log.Error().Err(err).Msg("error rendering template")
		return "Something went terribly wrong."
	}

	return style.Render(text)
}

var logger = log.With().Str("plugin", "templates").Logger()

var Plugin = &templatePlugin{
	themes:    make(map[string]theme.Theme),
	templates: make(map[string]template.Template),
}
