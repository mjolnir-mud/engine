package accounts

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/data_source"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
)

type plugin struct{}

func (p *plugin) Name() string {
	return "accounts"
}

func (p *plugin) Registered() error {
	engine.EnsureRegistered(data_sources.Plugin.Name())
	engine.EnsureRegistered(mongo_data_source.Plugin.Name())
	engine.EnsureRegistered(ecs.Plugin.Name())

	data_sources.Register(data_source.DataSource)
	return nil
}

func (p *plugin) Start() error {
	return nil
}

func (p *plugin) Stop() error {
	return nil
}

var Plugin = plugin{}
