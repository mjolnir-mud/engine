/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package entity

import (
	"context"

	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rueian/rueidis"
)

// Record is the record that is stored in Redis for an entity. It is stored as JSON using ReJSON.
type Record struct {
	// Id is the unique identifier for the entity.
	Id uid.UID

	// Version is the version of the entity. This is incremented every time a component is added, removed or updated.
	Version int64

	// Entity is the entity itself.
	Entity interface{}
}

//
//// GetComponent returns the named component for the given entity. If the entity or component does not exist, an error will
//// be returned. If the component is not found, an error will be returned.
//func GetComponent(ctx context.Context, entityOrId interface{}, componentName string, component interface{}) error {
//	entityId := mustGetEntityId(entityOrId)
//
//	logger := engine.GetLogger(ctx).
//		With().
//		Str("component", "entities").
//		Str("entityId", string(entityId)).
//		Str("componentName", componentName).
//		Logger()
//
//	logger.Debug().Msg("getting component")
//	logger.Trace().Msg("checking if entity exists")
//
//	logger.Trace().Msg("checking if component exists")
//	exists, err := HasComponent(ctx, entityId, componentName)
//
//	if err != nil {
//		return err
//	}
//
//	if !exists {
//		logger.Error().Msg("component does not exist")
//		return errors.ComponentNotFoundError{
//			EntityId: string(entityId),
//			Name:     componentName,
//		}
//	}
//
//	redis := engine.getRedisClient(ctx)
//
//	logger.Trace().Msg("building redis command")
//	command := redis.B().JsonGet().Key(string(entityId)).Paths(componentPath(componentName)).Build()
//	redis.UnmarshallResult(
//		redis.MustExecuteRedisCommands(ctx, rueidis.Commands{command})[0],
//		component,
//	)
//
//	if err != nil {
//		panic(err)
//	}
//
//	return nil
//}

//

//
//// HasComponent returns true if the component exists on the entity. If the entity does not exist an error is thrown.
//func HasComponent(ctx context.Context, entityOrId interface{}, componentName string) (bool, error) {
//	entityId, err := getEntityId(entityOrId)
//
//	if err != nil {
//		panic(err)
//	}
//
//	entityExists := HasEntity(ctx, entityId)
//
//	if !entityExists {
//		return false, errors.EntityNotFoundError{
//			Id: entityId,
//		}
//	}
//
//	result := redis.ExecuteRedisCommands(ctx, rueidis.Commands{
//		engine.getRedisClient(ctx).B().JsonType().Key(string(entityId)).Path(componentPath(componentName)).Build(),
//	})[0]
//
//	err = result.Error()
//
//	if err != nil {
//		if err.Error() == fmt.Sprintf("ERR Path '$%s' does not exist", componentPath(componentName)) {
//			return false, nil
//		}
//
//		return false, err
//	}
//
//	str, err := result.ToAny()
//
//	if err != nil {
//		panic(err)
//	}
//
//	if str == nil {
//		return false, nil
//	}
//
//	return true, nil
//}
//
//// RemoveEntity removes the entity from the engine. If the entity does not exist, an error will be returned. Each
//// component will be removed from the engine as well. This will trigger `events.EntityRemovedEvent` and a series of
//// `events.ComponentRemovedEvent` events for every component on the entity.
//func RemoveEntity(ctx context.Context, entity interface{}) error {
//	entityId := mustGetEntityId(entity)
//
//	logger := engine.GetLogger(ctx).With().Str("component", "entities").Str("entityId", string(entityId)).Logger()
//	logger.Debug().Msg("removing entity")
//
//	err := Get(ctx, entity)
//
//	if err != nil {
//		return err
//	}
//
//	redis.MustExecuteRedisCommandsAndPublishEvents(ctx,
//		rueidis.Commands{
//			engine.getRedisClient(ctx).B().Del().Key(string(entityId)).Build(),
//		},
//
//		buildEntityRemovedEvents(entityId, getComponentMap(entity)),
//	)
//
//	return nil
//}
//

//
//func buildEntityRemovedEvents(entityId uid.UID, components map[string]interface{}) []event.Event {
//	events := make([]event.Event, 0)
//
//	events = append(events, events2.EntityRemovedEvent{
//		Id: entityId,
//	})
//
//	events = append(events, buildComponentRemovedEvents(entityId, components)...)
//
//	return events
//}
//
//func buildComponentRemovedEvents(entityId uid.UID, components map[string]interface{}) []event.Event {
//	events := make([]event.Event, 0)
//
//	for name, value := range components {
//		events = append(events, buildComponentRemovedEvent(entityId, name, value))
//	}
//
//	return events
//}
//
//func buildComponentRemovedEvent(entityId uid.UID, componentName string, value interface{}) event.Event {
//	return events2.ComponentRemovedEvent{
//		EntityId: entityId,
//		Name:     componentName,
//		Value:    value,
//	}
//}
//

//

//	func componentPath(componentName string) string {
//		return fmt.Sprintf(".Entity.%s", componentName)
//	}
func getEntityRecord(ctx context.Context, id uid.UID) Record {
	var record Record

	redis.UnmarshallResult(
		redis.MustExecuteCommands(ctx, rueidis.Commands{
			redis.GetClient(ctx).B().JsonGet().Key(string(id)).Build(),
		})[0],
		&record,
	)

	return record
}
