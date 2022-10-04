package plugin

// Plugin is the data_source that must be implemented by a Mjolnir plugin. The game is expected to call
// `engine.RegisterPlugin(plugin)` for every plugin that is to be used for that game before the `engine.Start()`
// is called. See [the plugins wiki documentation](https://github.com/mjolnir-mud/engine/wiki/Plugins) for more
// information.
type Plugin interface {
	// Name returns the name of the plugin. Plugin names must be unique.
	Name() string

	// Registered is a callback that is called when the plugin is registered with the game. Developers can utilize this
	// callback to perform any initialization that is required before the game engine is started.
	Registered() error
}