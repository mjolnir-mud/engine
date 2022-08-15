package template_registry

import (
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/template"
	"github.com/rs/zerolog"
)

var templates = map[string]template.Template{}
var log zerolog.Logger

func Start() {
	log = logger.Instance.With().Str("service", "template_registry").Logger()
}

func Stop() {}

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

	s, err := t.Render(data)

	if err != nil {
		return "", err
	}

	return s, nil
}
