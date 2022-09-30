package helpers

import "github.com/mjolnir-mud/engine/plugins/data_sources"

func DeleteAccount(id string) {
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{
		"id": id,
	})
}

func DeleteDefaultAccount() {
	DeleteAccount(TestAccountId)
}
