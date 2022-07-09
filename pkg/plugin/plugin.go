package plugin

// Plugin is the interface that must be implemented by a Mjolnir plugin.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string

	// Init initializes the plugin when the game starts.
	Init() error
}
