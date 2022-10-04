package theme_registry

import (
	"github.com/mjolnir-mud/engine/plugins/templates/errors"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/templates/theme"
	"github.com/rs/zerolog"
)

var themes map[string]theme.Theme
var log zerolog.Logger
var globalTheme = "default"

func Start() {
	log = logger.Instance.With().Str("component", "theme_registry").Logger()
	themes = map[string]theme.Theme{}
	log.Info().Msg("starting theme registry")
}

func Stop() {
	log.Info().Msg("stopping theme registry")
}

func Register(t theme.Theme) {
	log.Info().Str("theme", t.Name()).Msg("registering theme")
	themes[t.Name()] = t
}

func GetTheme(name string) (theme.Theme, error) {
	t, ok := themes[name]

	if !ok {
		return nil, errors.ThemeNotFoundError{
			Name: name,
		}
	}

	return t, nil
}

func Render(style string, content string) (string, error) {
	t, err := GetTheme(globalTheme)

	if err != nil {
		return "", err
	}

	s := t.GetStyleFor(style)

	if err != nil {
		return "", err
	}

	return s.Render(content), nil
}

func SetGlobalTheme(name string) {
	globalTheme = name
}
