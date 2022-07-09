package engine

type testPlugin struct {
	pluginInitialized chan bool
}

func (p *testPlugin) Name() string {
	return "test"
}

func (p *testPlugin) Init() error {
	p.pluginInitialized <- true

	return nil
}

func createTestPlugin() Plugin {
	return &testPlugin{
		pluginInitialized: make(chan bool),
	}
}

//func TestLoadPlugins(t *testing.T) {
//	Init("test", []Plugin{&testPlugin{}})
//
//	assert.True(t, initCalled)
//
//	assert.NotNilf(t, state.baseCommand, "baseCommand should not be nil")
//}
