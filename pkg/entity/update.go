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
	"reflect"

	"github.com/mjolnir-engine/engine/pkg/errors"
	"github.com/mjolnir-engine/engine/pkg/events"
	"github.com/mjolnir-engine/engine/pkg/logger"
	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rueian/rueidis"
)

// Update updates a component on an entity. This will trigger the `events.ComponentUpdatedEvent` event to be
// published for every component that has changed, if that component already exists on the entity. It will trigger the
// `events.ComponentAddedEvent` event for every component that did not exist on the entity. It will trigger the
// `events.ComponentRemovedEvent` for any components that no longer exist on the entity. It will trigger the
// `events.EntityUpdatedEvent` event if any component has changed. If the entity does not exist, an error will be
// returned.
func Update(ctx context.Context, entity interface{}) error {
	err := typeCheckStructPointer(entity)

	if err != nil {
		return err
	}

	m := getComponentMap(entity)

	var id uid.UID
	id, ok := m["Id"].(uid.UID)

	if !ok {
		return errors.EntityNotFoundError{
			Id: m["Id"].(uid.UID),
		}
	}

	l := logger.Get(ctx).With().
		Str("component", "entities").
		Str("entityId", string(id)).
		Logger()

	l.Debug().Msg("updating component")

	l.Trace().Msg("checking if entity exists")

	if !HasEntity(ctx, id) {
		return errors.EntityNotFoundError{
			Id: id,
		}
	}

	entityRecord := getEntityRecord(ctx, id)

	oldEntityMap := entityRecord.Entity.(map[string]interface{})
	newEntityMap := getComponentMap(entity)

	if err != nil {
		return err
	}

	l.Trace().Msg("comparing old and new entity")

	if reflect.DeepEqual(oldEntityMap, newEntityMap) {
		l.Trace().Msg("entity has not changed")
		return nil
	}

	entityRecord.Version++
	entityRecord.Entity = entity

	// set the updated entity
	commands := rueidis.Commands{
		redis.GetClient(ctx).B().JsonSet().Key(string(id)).Path(".").Value(rueidis.JSON(entityRecord)).Build(),
	}

	l.Trace().Msg("building publish commands for events")

	l.Trace().Msg("executing redis commands")

	redis.MustExecuteRedisCommandsAndPublishEvents(
		ctx,
		commands,
		buildComponentUpdatedEvents(newEntityMap, oldEntityMap),
	)

	return nil
}

func buildComponentUpdatedEvents(entity map[string]interface{}, oldEntity map[string]interface{}) []events.Event {
	id, ok := entity["Id"].(uid.UID)

	if !ok {
		panic("unable to get entity id")
	}

	e := make([]events.Event, 0)

	for name, value := range entity {
		if name == "Id" {
			continue
		}

		if oldEntity[name] == nil {
			e = append(e, events.ComponentAddedEvent{
				EntityId: id,
				Name:     name,
				Value:    value,
			})
		}

		if !reflect.DeepEqual(value, oldEntity[name]) {
			e = append(e, events.ComponentUpdatedEvent{
				EntityId:      id,
				Name:          name,
				Value:         value,
				PreviousValue: oldEntity[name],
			})
		}
	}

	for name, value := range oldEntity {
		if entity[name] == nil {
			e = append(e, events.ComponentRemovedEvent{
				EntityId: id,
				Name:     name,
				Value:    value,
			})
		}
	}

	return e
}
