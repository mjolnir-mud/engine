package theme_registry

import (
	"github.com/mjolnir-mud/engine/plugins/templates/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/templates/pkg/theme"
	"github.com/rs/zerolog"
)

var themes map[string]theme.Theme
var log zerolog.Logger

func Start() {
	log = logger.Instance.With().Str("service", "theme_registry").Logger()
	themes = map[string]theme.Theme{}
}

func Stop() {}

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
