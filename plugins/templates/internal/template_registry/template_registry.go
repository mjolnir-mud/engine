package template_registry

import (
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/theme_registry"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/template"
	"github.com/rs/zerolog"
)

var templates map[string]template.Template
var log zerolog.Logger

func Start() {
	templates = make(map[string]template.Template)
	log = logger.Instance.With().Str("service", "template_registry").Logger()
	log.Info().Msg("starting template registry")
}

func Stop() {
	log.Info().Msg("stopping template registry")
}

func Register(t template.Template) {
	log.Info().Str("template", t.Name()).Msg("registering template")
	templates[t.Name()] = t
}

func Render(name string, data interface{}) (string, error) {
	t, ok := templates[name]

	if !ok {
		return "", errors.TemplateNotFoundError{
			Name: name,
		}
	}

	th, err := theme_registry.GetTheme("default")

	if err != nil {
		return "", err
	}

	style := th.GetStyleFor(t.Style())
	text, err := t.Render(data)

	if err != nil {
		return "", err
	}

	return style.Render(text), nil
}
