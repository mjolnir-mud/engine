package engine

type pluginManager struct {
	plugins []Plugin
}

func newPluginManager() *pluginManager {
	return &pluginManager{}
}
