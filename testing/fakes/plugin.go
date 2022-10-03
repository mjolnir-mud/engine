package fakes

import "github.com/mjolnir-mud/engine/pkg/plugin"

func CreateFakePlugin() plugin.Plugin {
	return &testPlugin{
		pluginInitialized: make(chan bool),
	}
}

type testPlugin struct {
	pluginInitialized chan bool
}

func (p *testPlugin) Name() string {
	return "testing"
}

func (p *testPlugin) Start() error {
	p.pluginInitialized <- true

	return nil
}

func (p *testPlugin) Registered() error {
	return nil
}

func (p *testPlugin) Stop() error {
	return nil
}
