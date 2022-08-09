package accounts

import (
	"github.com/mjolnir-mud/engine/plugins/accounts/pkg/data_source"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
)

type plugin struct{}

func (p *plugin) Name() string {
	return "accounts"
}

func (p *plugin) Registered() error {
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
