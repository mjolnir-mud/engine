package events

import (
	"fmt"

	"github.com/mjolnir-mud/engine/pkg/engine"
)

type PluginReady struct {
	Name string
}

// PluginReadyTopic returns the topic for the plugin ready event.
func PluginReadyTopic(plugin engine.Plugin) string {
	return fmt.Sprintf("%s.ready", plugin.Name())
}
