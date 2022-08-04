package test

import "github.com/mjolnir-mud/engine/pkg/plugin"

func CreateTestPlugin() plugin.Plugin {
	return &testPlugin{
		pluginInitialized: make(chan bool),
	}
}

type testPlugin struct {
	pluginInitialized chan bool
}

func (p *testPlugin) Name() string {
	return "test"
}

func (p *testPlugin) Start() error {
	p.pluginInitialized <- true

	return nil
}
