package engine

// Plugin is the interface that must be implemented by a Mjolnir plugin.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string

	// Init initializes the plugin when the game starts.
	Init(state *State) error
}

// LoadPlugins loads all the plugins in the game state, passing the state to each plugin's `Init` function.
func LoadPlugins(state *State) error {
	for _, plugin := range state.Plugins {
		if err := plugin.Init(state); err != nil {
			return err
		}
	}
	return nil
}
