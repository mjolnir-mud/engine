package data_source

import "github.com/mjolnir-mud/engine/plugins/mongo_data_source"

var DataSource = mongo_data_source.New("accounts")