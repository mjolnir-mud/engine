package helpers

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	sessionsTesting "github.com/mjolnir-mud/engine/plugins/sessions/testing/helpers"
)

const TestAccountId = "test-account"

func CreateSessionWithOutputSubscription() (string, chan string, engine.Subscription, error) {
	id, ch, sub, err := sessionsTesting.CreateSessionWithOutputSubscription()

	if err != nil {
		return "", nil, nil, err
	}

	err = ecs.AddStringComponentToEntity(id, "accountId", TestAccountId)

	if err != nil {
		return "", nil, nil, err
	}

	return id, ch, sub, nil
}

func CreateSessionWithOutputSubscriptionForAccountId(accountId string) (string, chan string, engine.Subscription, error) {
	id, ch, sub, err := sessionsTesting.CreateSessionWithOutputSubscription()

	if err != nil {
		return "", nil, nil, err
	}

	err = ecs.AddStringComponentToEntity(id, "accountId", accountId)

	if err != nil {
		return "", nil, nil, err
	}

	return id, ch, sub, nil
}
