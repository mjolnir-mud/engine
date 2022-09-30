package helpers

import "github.com/mjolnir-mud/engine/plugins/data_sources"

func CreateAccount(id string) {
	_, _ = data_sources.CreateEntityWithId(
		"accounts",
		"account",
		id,
		map[string]interface{}{},
	)
}

func CreateDefaultAccount() {
	CreateAccount(TestAccountId)
}
