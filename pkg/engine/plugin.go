package engine

// Plugin is the interface that must be implemented by a Mjolnir plugin.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string

	// Init initializes the plugin when the game starts.
	Init() error
}

func loadPlugins() error {
	for _, plugin := range state.plugins {
		if err := plugin.Init(); err != nil {
			return err
		}
	}
	return nil
}
