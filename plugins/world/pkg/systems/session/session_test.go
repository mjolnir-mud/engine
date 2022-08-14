package session

import (
	"github.com/mjolnir-mud/engine/plugins/ecs"
	session2 "github.com/mjolnir-mud/engine/plugins/world/pkg/entities/session"
	testing2 "github.com/mjolnir-mud/engine/plugins/world/pkg/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	ecs.RegisterEntityType(session2.Type)
	testing2.Setup()

	ent, err := ecs.CreateEntity("test", map[string]interface{}{})

	if err != nil {
		panic(err)
	}

	err = ecs.AddEntityWithID("test", "test", ent)

	if err != nil {
		panic(err)
	}
}

func teardown() {
	testing2.Teardown()
}

func TestStart(t *testing.T) {
	testing2.Setup()
	defer testing2.Teardown()

	err := Start("test")

	assert.NoError(t, err)
}
