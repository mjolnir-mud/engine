package account

import (
	"github.com/mjolnir-mud/engine/plugins/data_sources/data_source"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
)

func Create() data_source.Interface {
	return mongo_data_source.New("accounts")
}
