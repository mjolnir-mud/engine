package data_source

import (
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/data_source"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
)

func Create() data_source.DataSource {
	return mongo_data_source.New("accounts")
}
