package metadata

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/constants"
)

func SetMetadata(entityId string, key string, value interface{}) error {
	err := setTypeMetadata(entityId, key, value)

	if err != nil {
		return err
	}

	return nil
}

func RemoveMetadata(entityId string, key string) error {
	return nil
}

func HasMetadata(entityId string, key string, metadataKey string) bool {
	engine.Redis.Exists(context.Background(), fmt.Sprintf("__%s:%s:%s", metadataKey, entityId))
}

func setTypeMetadata(entityId string, key string, value interface{}) error {
	if !HasMetadata(entityId, key, constants.ComponentTypePrefix) {
		valueType := reflect.TypeOf(value).Kind().String()

		switch valueType {
		case "slice":
			valueType = "set"
		}

		return engine.Redis.Set(
			context.Background(),
			fmt.Sprintf("__%s:%s:%s", constants.ComponentTypePrefix, entityId, key),
			valueType,
			0,
		).Err()
	}

	return nil
}
