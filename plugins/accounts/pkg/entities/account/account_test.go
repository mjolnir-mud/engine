package account

import (
	"github.com/mjolnir-mud/engine"
	testing2 "github.com/mjolnir-mud/engine/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/accounts/internal/data_source"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func setup() {
	engine.RegisterPlugin(data_sources.Plugin)
	engine.RegisterPlugin(mongo_data_source.Plugin)
	engine.RegisterPlugin(ecs.Plugin)
	ecs.RegisterEntityType(Type)
	data_sources.Register(data_source.Create())
	testing2.Setup()

	// seed an account
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testaccount"})

	err := data_sources.Save(
		"accounts",
		"testaccount",
		map[string]interface{}{
			"username": "testaccount",
			"password": string(hashedPassword),
			"__metadata": map[string]interface{}{
				"entityType": "account",
				"collection": "accounts",
			},
		})

	if err != nil {
		panic(err)
	}
}

func teardown() {
	_ = data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testaccount"})
	testing2.Teardown()
}

func TestAccountType_Name(t *testing.T) {
	setup()
	defer teardown()
	assert.Equal(t, "account", Type.Name())
}

func TestAccountType_Create(t *testing.T) {
	assert.Equal(t, map[string]interface{}{"username": "test"}, Type.Create(map[string]interface{}{"username": "test"}))
}

func TestValidateAccount(t *testing.T) {
	setup()
	defer teardown()

	id, err := ValidateAccount(Credentials{"testaccount", "password"})

	assert.Nil(t, err)
	assert.Equal(t, "testaccount", id)

	_, err = ValidateAccount(Credentials{"testaccount", "wrongpassword"})

	assert.NotNil(t, err)
}
