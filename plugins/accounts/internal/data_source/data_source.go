package data_source

import (
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
)

func Create() *mongo_data_source.MongoDataSource {
	return mongo_data_source.New("accounts")
}
