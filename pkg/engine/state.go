package engine

type State struct {
	Plugins []Plugin
}

func newState(plugins []Plugin) *State {
	return &State{
		plugins,
	}
}
