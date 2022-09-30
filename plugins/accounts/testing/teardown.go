package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/accounts/testing/helpers"
)

func Teardown() {
	_ = engine.RedisFlushAll()
	helpers.DeleteDefaultAccount()
}
