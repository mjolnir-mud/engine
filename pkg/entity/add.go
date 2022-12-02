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

	"github.com/mitchellh/mapstructure"
	"github.com/mjolnir-engine/engine/pkg/errors"
	"github.com/mjolnir-engine/engine/pkg/event"
	"github.com/mjolnir-engine/engine/pkg/events"
	"github.com/mjolnir-engine/engine/pkg/logger"
	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rueian/rueidis"
)

// Add adds an entity to the engine. If the entity id is not set, it will be set to a new UUID. If the entity
// already exists, an error will be returned. This will trigger `events.EntityAddedEvent` and a `events.ComponentAddedEvent`
// for each component on the entity.
func Add(ctx context.Context, entity interface{}) error {
	err := typeCheckStructPointer(entity)

	if err != nil {
		return err
	}

	l := logger.Get(ctx).With().Str("component", "entities").Logger()
	l.Debug().Msg("adding entity")
	l.Trace().Msg("prepping id")

	id, err := getEntityId(entity)

	if err != nil || id == "" {
		id = uid.New()
		m := getComponentMap(entity)
		m["Id"] = id

		err := mapstructure.Decode(m, entity)

		if err != nil {
			panic(err)
		}
	}

	l = l.With().Str("entityId", string(id)).Logger()
	l.Trace().Msg("checking if entity exists")
	exists := HasEntity(ctx, entity)

	if err != nil {
		return err
	}

	if exists {
		return errors.EntityExistsError{
			Id: id,
		}
	}

	l.Trace().Msg("building client command")

	client := redis.GetClient(ctx)

	record := Record{
		Id:      string(id),
		Version: 1,
		Entity:  entity,
	}

	commands := rueidis.Commands{
		client.B().JsonSet().Key(string(id)).Path(".").Value(rueidis.JSON(record)).Build(),
	}

	components := getComponentMap(entity)

	redis.MustExecuteRedisCommandsAndPublishEvents(
		ctx,
		commands,
		buildEntityAndComponentAddedEvents(id, components),
	)

	l.Trace().Msg("entity added")
	return nil
}

func buildEntityAddedEvent(entityId uid.UID) event.Event {
	return events.EntityAddedEvent{
		Id: entityId,
	}
}

func buildEntityAndComponentAddedEvents(entityId uid.UID, components map[string]interface{}) []event.Event {
	events := []event.Event{
		buildEntityAddedEvent(entityId),
	}
	return append(events, buildComponentAddedEvents(entityId, components)...)
}
