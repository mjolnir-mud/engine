package helpers

import (
	"github.com/mjolnir-mud/engine/plugins/ecs"
	sessionsHelpers "github.com/mjolnir-mud/engine/plugins/sessions/testing/helpers"
)

func CreateSession() (string, error) {
	id, err := sessionsHelpers.CreateSession()

	if err != nil {
		return "", err
	}

	err = ecs.AddStringComponentToEntity(id, "accountId", TestAccountId)

	if err != nil {
		return "", err
	}

	return id, nil
}
