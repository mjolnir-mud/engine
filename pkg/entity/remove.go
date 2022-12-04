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

	"github.com/mjolnir-engine/engine/pkg/events"
	"github.com/mjolnir-engine/engine/pkg/logger"
	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rueian/rueidis"
)

// Remove removes the entity from the engine. If the entity does not exist, an error will be returned. Each
// component will be removed from the engine as well. This will trigger `events.EntityRemovedEvent` and a series of
// `events.ComponentRemovedEvent` events for every component on the entity.
func Remove(ctx context.Context, entity interface{}) error {
	entityId := mustGetEntityId(entity)

	l := logger.Get(ctx).With().Str("component", "entities").Str("entityId", string(entityId)).Logger()
	l.Debug().Msg("removing entity")

	err := Get(ctx, entity)

	if err != nil {
		return err
	}

	redis.MustExecuteRedisCommandsAndPublishEvents(ctx,
		rueidis.Commands{
			redis.GetClient(ctx).B().Del().Key(string(entityId)).Build(),
		},

		buildEntityRemovedEvents(entityId, getComponentMap(entity)),
	)

	return nil
}

func buildEntityRemovedEvents(entityId uid.UID, components map[string]interface{}) []events.Event {
	e := make([]events.Event, 0)

	e = append(e, events.EntityRemovedEvent{
		Id: entityId,
	})

	e = append(e, buildComponentRemovedEvents(entityId, components)...)

	return e
}

func buildComponentRemovedEvents(entityId uid.UID, components map[string]interface{}) []events.Event {
	evs := make([]events.Event, 0)

	for name, value := range components {
		evs = append(evs, buildComponentRemovedEvent(entityId, name, value))
	}

	return evs
}

func buildComponentRemovedEvent(entityId uid.UID, componentName string, value interface{}) events.Event {
	return events.ComponentRemovedEvent{
		EntityId: entityId,
		Name:     componentName,
		Value:    value,
	}
}
