package account

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/accounts/data_sources/account"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	testing2 "github.com/mjolnir-mud/engine/plugins/data_sources/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	testing3 "github.com/mjolnir-mud/engine/plugins/mongo_data_source/testing"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func setup() {
	engineTesting.RegisterSetupCallback("accounts", func() {
		testing2.Setup()
		ecsTesting.Setup()
		testing3.Setup()

		engine.RegisterBeforeServiceStartCallback("world", func() {
			data_sources.Register(account.Create())
		})

		engine.RegisterAfterServiceStartCallback("world", func() {
			ecs.RegisterEntityType(EntityType)

			_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testaccount"})
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			err := data_sources.SaveWithId(
				"accounts",
				"testaccount",
				map[string]interface{}{
					"username":       "testaccount",
					"hashedPassword": string(hashedPassword),
					"__metadata": map[string]interface{}{
						"entityType": "account",
						"collection": "accounts",
					},
				})

			if err != nil {
				panic(err)
			}
		})
	})

	engineTesting.Setup("world")
}

func teardown() {
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testaccount"})
	testing2.Teardown()
	testing3.Teardown()
	engineTesting.Teardown()
}

func TestAccountType_Name(t *testing.T) {
	setup()
	defer teardown()
	assert.Equal(t, "account", EntityType.Name())
}

func TestAccountType_New(t *testing.T) {
	assert.Equal(t, map[string]interface{}{"username": "testing"}, EntityType.New(map[string]interface{}{"username": "testing"}))
}

func TestAccountType_Validate(t *testing.T) {
	setup()
	defer teardown()

	assert.Error(t, EntityType.Validate(map[string]interface{}{"username": "testing"}))
}
