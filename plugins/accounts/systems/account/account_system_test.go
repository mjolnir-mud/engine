package account

import (
	"github.com/mjolnir-mud/engine"
	engineTesting "github.com/mjolnir-mud/engine/pkg/testing"
	accountDataSource "github.com/mjolnir-mud/engine/plugins/accounts/data_sources/account"
	accountEntityType "github.com/mjolnir-mud/engine/plugins/accounts/entities/account"
	"github.com/mjolnir-mud/engine/plugins/data_sources"
	dataSourcesTesting "github.com/mjolnir-mud/engine/plugins/data_sources/testing"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	mongoDataSourceTesting "github.com/mjolnir-mud/engine/plugins/mongo_data_source/testing"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	engineTesting.Setup("world", func() {
		dataSourcesTesting.Setup()
		ecsTesting.Setup()
		mongoDataSourceTesting.Setup()

		engine.RegisterBeforeServiceStartCallback("world", func() {
			data_sources.Register(accountDataSource.Create())
		})

		engine.RegisterAfterServiceStartCallback("world", func() {
			ecs.RegisterEntityType(accountEntityType.EntityType)

			waitForResult := make(chan interface{})

			go func() {
				waitForResult <- data_sources.FindAndDelete("accounts", map[string]interface{}{"username": "testing-account"})
			}()

			<-waitForResult

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			err := data_sources.SaveWithId(
				"accounts",
				"testing-account",
				map[string]interface{}{
					"username":       "testing-account",
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
}

func teardown() {
	ecsTesting.Teardown()
	engineTesting.Teardown()
}

func TestSystem_Name(t *testing.T) {
	setup()
	defer teardown()
	assert.Equal(t, "account", System.Name())
}

func TestSystem_Component(t *testing.T) {
	setup()
	defer teardown()
	assert.Equal(t, "_", System.Component())
}

func TestValidateEmail(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ValidateEmail("testing@email.com"))

	assert.Error(t, ValidateEmail("testing"))
}

func TestValidatePassword(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ValidatePassword("testing", "testing@testing.com", "validPassword"))
	assert.Error(t, ValidatePassword("testing", "testing@testing.com", "invalid"))
	assert.Error(t, ValidatePassword("validPassword", "testing@testing.com", "validPassword"))
	assert.Error(t, ValidatePassword("testing", "testing@testing.com", "testing@testing.com"))
}

func TestValidateUsername(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, ValidateUsername("testing"))
	assert.Error(t, ValidateUsername("invalid@"))
}

func TestCompareAccountCredentials(t *testing.T) {
	setup()
	defer teardown()

	id, err := CompareAccountCredentials(Credentials{"testing-account", "password"})

	assert.Nil(t, err)
	assert.NotNil(t, id)

	_, err = CompareAccountCredentials(Credentials{"testing-account", "wrongpassword"})

	assert.NotNil(t, err)
}
