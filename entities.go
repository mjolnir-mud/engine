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

package engine

import (
	"context"
	"fmt"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	engineErrors "github.com/mjolnir-mud/engine/errors"
	engineEvents "github.com/mjolnir-mud/engine/events"
	"github.com/mjolnir-mud/engine/internal/uid"
	"github.com/rueian/rueidis"
	"reflect"
)

// EntityRecord is the record that is stored in Redis for an entity. It is stored as JSON using ReJSON.
type EntityRecord struct {
	// Id is the unique identifier for the entity.
	Id string

	// Version is the version of the entity. This is incremented every time a component is added, removed or updated.
	Version int64

	// Entity is the entity itself.
	Entity interface{}
}

// AddComponent Adds a component to an entity. This will trigger the `events.ComponentAddedEvent` event to be published.
// If the entity does not exist, an error will be returned. If the component already exists, an error will be returned.
// If you wish to update a component, use the `UpdateComponent` method.
func (e *Engine) AddComponent(entityId string, componentName string, component interface{}) error {
	return nil
}

// AddEntity adds an entity to the engine with a random id, returning the id. This will trigger the `events.EntityAddedEvent`
// event to be published. Then entity must be a struct, otherwise an error will be returned.
func (e *Engine) AddEntity(entity interface{}) (string, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	return id.String(), e.AddEntityWithId(id.String(), entity)
}

// AddEntityWithId adds an entity to the engine. This will trigger the `events.EntityAddedEvent` event to be published. If
// the entity already exists, an error will be returned. The id will be converted to a Mjolnir UID before it is added.
// Then entity must be a struct, otherwise an error will be returned.
func (e *Engine) AddEntityWithId(id string, entity interface{}) error {
	logger := e.Logger.With().Str("component", "entities").Str("entityId", id).Logger()

	if !entityIsStruct(entity) {
		return engineErrors.EntityInvalidError{Id: id, Value: entity}
	}

	logger.Trace().Msg("checking if entity exists")
	exists, err := e.HasEntity(id)

	if err != nil {
		return err
	}

	if exists {
		logger.Error().Msg("entity already exists")
		return engineErrors.EntityExistsError{
			Id: id,
		}
	}

	eKey := e.stringToKey(id)

	record := EntityRecord{
		Id:      eKey,
		Version: 1,
		Entity:  entity,
	}

	logger.Trace().Msg("building redis commands")
	commands := rueidis.Commands{
		e.redis.B().JsonSet().Key(eKey).Path(".").Value(rueidis.JSON(record)).Build(),
	}

	componentNames := getComponentNames(entity)
	logger.Trace().Strs("components", componentNames).Msg("building publish commands for components")

	commands = append(commands, e.GetPublishCommandsForEvents(buildEntityAndComponentAddedEvents(id, componentNames)...)...)

	logger.Trace().Msg("executing redis commands")
	results := e.redis.DoMulti(
		context.Background(),
		commands...,
	)

	logger.Trace().Msg("checking redis results")
	for _, result := range results {
		if result.Error() != nil {
			logger.Error().Err(result.Error()).Msg("error adding entity")
		}
	}

	logger.Trace().Msg("component added")

	return nil
}

// FlushEntities removes all entities from the engine.
func (e *Engine) FlushEntities() error {
	keys, err := e.redis.Do(context.Background(), e.redis.B().Keys().Pattern(e.uidToKey("*")).Build()).AsStrSlice()

	for _, key := range keys {
		_, err = e.redis.Do(context.Background(), e.redis.B().Del().Key(key).Build()).AsBool()

		if err != nil {
			return err
		}
	}

	return nil
}

// HasEntity returns true if the entity exists in the engine.
func (e *Engine) HasEntity(id string) (bool, error) {
	exists, err := e.redis.Do(context.Background(), e.redis.B().Exists().Key(e.stringToKey(id)).Build()).AsBool()

	if err != nil {
		return false, err
	}

	return exists, nil
}

// HasComponent returns true if the component exists on the entity.
func (e *Engine) HasComponent(entityId string, componentName string) (bool, error) {
	res, err := e.redis.Do(
		context.Background(),
		e.redis.B().JsonGet().Key(e.stringToKey(entityId)).Paths(componentPath(componentName)).Build(),
	).ToMessage()

	if err != nil {
		if err.Error() == fmt.Sprintf("ERR Path '$%s' does not exist", componentPath(componentName)) {
			return false, nil
		}

		return false, err
	}

	str, err := res.ToAny()

	if err != nil {
		return false, err
	}

	if str == nil {
		return false, nil
	}

	return true, nil
}

// UpdateComponent updates a component on an entity. This will trigger the `events.ComponentUpdatedEvent` event to be
// published. If the entity does not exist, an error will be returned. If the component does not exist, an error will be
// returned.
func (e *Engine) UpdateComponent(entityId string, componentName string, component interface{}) error {
	return nil
}

func (e *Engine) uidToKey(id string) string {
	return fmt.Sprintf("%s:entity:%s", e.instanceId, id)
}

func (e *Engine) stringToKey(id string) string {
	return e.uidToKey(uid.FromString(id))
}

func (e *Engine) publishComponentsAdded(entityId string, componentNames []string) error {
	return e.Publish(buildComponentAddedEvents(entityId, componentNames)...)
}

func (e *Engine) publishEntityAndComponentsAdded(entityId string, componentNames []string) error {
	return e.Publish(buildEntityAndComponentAddedEvents(entityId, componentNames)...)
}

func buildComponentAddedEvents(entityId string, componentNames []string) []Event {
	events := make([]Event, len(componentNames))

	for i, componentName := range componentNames {
		events[i] = engineEvents.ComponentAddedEvent{
			EntityId: entityId,
			Name:     componentName,
		}
	}

	return events
}

func buildEntityAddedEvent(entityId string) Event {
	return engineEvents.EntityAddedEvent{
		Id: entityId,
	}
}

func buildEntityAndComponentAddedEvents(entityId string, componentNames []string) []Event {
	events := []Event{
		buildEntityAddedEvent(entityId),
	}
	return append(events, buildComponentAddedEvents(entityId, componentNames)...)
}

func getComponentNames(entity interface{}) []string {
	return structs.Names(entity)
}

func componentPath(componentName string) string {
	return fmt.Sprintf(".Entity.%s", componentName)
}

func entityIsStruct(entity interface{}) bool {
	if entity == nil {
		return false
	}

	return reflect.TypeOf(entity).Kind() == reflect.Struct
}
