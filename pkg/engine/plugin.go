package engine

// Plugin is the interface that must be implemented by a Mjolnir plugin.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string

	// Init initializes the plugin when the game starts.
	Init(state *State) error
}

func loadPlugins(state *State) error {
	for _, plugin := range state.Plugins {
		if err := plugin.Init(state); err != nil {
			return err
		}
	}
	return nil
}
