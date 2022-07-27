package plugin

// Plugin is the interface that must be implemented by a Mjolnir plugin. The game is expected to call
// `engine.RegisterPlugin(plugin)` for every plugin that is to be used for that game before the `engine.Start()`
// is called. See [the plugins wiki documentation](https://github.com/mjolnir-mud/engine/wiki/Plugins) for more
// information.
type Plugin interface {
	// Name returns the name of the plugin. Plugin names must be unique.
	Name() string

	// Start initializes the plugin when the game starts. This is where any code that the plugin wishes to execute
	// during the startup process should be defined.
	Start() error
}
