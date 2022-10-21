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
	engineErrors "github.com/mjolnir-engine/engine/errors"
	engineEvents "github.com/mjolnir-engine/engine/events"
	"github.com/mjolnir-engine/engine/uid"
	"github.com/rueian/rueidis"
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
func (e *Engine) AddComponent(entityId uid.UID, componentName string, component interface{}) error {
	logger := e.logger.
		With().
		Str("component", "entities").
		Str("entityId", string(entityId)).
		Str("componentName", componentName).
		Logger()

	logger.Debug().Msg("adding component")

	logger.Trace().Msg("checking if entity exists")
	exists, err := e.HasEntity(entityId)

	if err != nil {
		return err
	}

	if !exists {
		logger.Error().Msg("entity does not exist")
		return engineErrors.EntityNotFoundError{
			Id: entityId,
		}
	}

	logger.Trace().Msg("checking if component exists")
	exists, err = e.HasComponent(entityId, componentName)

	if err != nil {
		return err
	}

	if exists {
		logger.Error().Msg("component already exists")
		return engineErrors.ComponentExistsError{
			EntityId: string(entityId),
			Name:     componentName,
		}
	}

	logger.Trace().Msg("building redis command")
	commands := rueidis.Commands{
		e.redis.
			B().
			JsonSet().
			Key(string(entityId)).Path(componentPath(componentName)).
			Value(rueidis.JSON(component)).
			Build(),
		e.redis.
			B().
			JsonNumincrby().
			Key(string(entityId)).
			Path(".Version").
			Value(1).
			Build(),
	}

	componentMap := map[string]interface{}{
		componentName: component,
	}

	commands = append(commands, e.GetPublishCommandsForEvents(buildComponentAddedEvents(entityId, componentMap)...)...)

	logger.Trace().Msg("executing redis command")
	results := e.redis.DoMulti(context.Background(), commands...)

	logger.Trace().Msg("checking redis results")

	cae := engineErrors.AddComponentErrors{}
	for _, result := range results {
		if result.Error() != nil {
			cae.Add(result.Error())
		}
	}

	if cae.HasErrors() {
		return cae
	}

	return nil
}

// AddEntity adds an entity to the engine with a random id, returning the id. This will trigger the
// `events.EntityAddedEvent` event to be published. Then entity must be a struct, otherwise an error will be returned.
func (e *Engine) AddEntity(entity interface{}) (uid.UID, error) {
	id := uid.New()

	return id, e.AddEntityWithId(id, entity)
}

// AddEntityWithId adds an entity to the engine. This will trigger the `events.EntityAddedEvent` event to be published.
// If the entity already exists, an error will be returned. The id will be converted to a Mjolnir UID before it is
// added. Then entity must be a struct, otherwise an error will be returned.
func (e *Engine) AddEntityWithId(id uid.UID, entity interface{}) error {
	logger := e.logger.With().Str("component", "entities").Str("entityId", string(id)).Logger()
	logger.Debug().Msg("adding entity")

	if !entityIsStruct(entity) {
		return engineErrors.EntityInvalidError{Id: string(id), Value: entity}
	}

	logger.Trace().Msg("checking if entity exists")
	exists, err := e.HasEntity(id)

	if err != nil {
		return err
	}

	if exists {
		logger.Error().Msg("entity already exists")
		return engineErrors.EntityExistsError{
			Id: string(id),
		}
	}

	record := EntityRecord{
		Id:      string(id),
		Version: 1,
		Entity:  entity,
	}

	logger.Trace().Msg("building redis commands")
	commands := rueidis.Commands{
		e.redis.B().JsonSet().Key(string(id)).Path(".").Value(rueidis.JSON(record)).Build(),
	}

	components := getComponentMap(entity)
	logger.Trace().Interface("components", components).Msg("building publish commands for components")

	commands = append(commands, e.GetPublishCommandsForEvents(
		buildEntityAndComponentAddedEvents(id, components)...)...,
	)

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
// TODO: This should actually do something.
func (e *Engine) FlushEntities() error {
	//keys, err := e.redis.Do(context.Background(), e.redis.B().Keys().Pattern(e.uidToKey("*")).Build()).AsStrSlice()
	//
	//for _, key := range keys {
	//	_, err = e.redis.Do(context.Background(), e.redis.B().Del().Key(key).Build()).AsBool()
	//
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

// GetComponent returns the named component for the given entity. If the entity or component does not exist, an error will
// be returned. If the component is not found, an error will be returned.
func (e *Engine) GetComponent(entityId uid.UID, componentName string, component interface{}) error {
	logger := e.logger.
		With().
		Str("component", "entities").
		Str("entityId", string(entityId)).
		Str("componentName", componentName).
		Logger()

	logger.Debug().Msg("getting component")

	logger.Trace().Msg("checking if entity exists")
	exists, err := e.HasEntity(entityId)

	if err != nil {
		return err
	}

	if !exists {
		logger.Error().Msg("entity does not exist")
		return engineErrors.EntityNotFoundError{
			Id: entityId,
		}
	}

	logger.Trace().Msg("checking if component exists")
	exists, err = e.HasComponent(entityId, componentName)

	if err != nil {
		return err
	}

	if !exists {
		logger.Error().Msg("component does not exist")
		return engineErrors.ComponentNotFoundError{
			EntityId: string(entityId),
			Name:     componentName,
		}
	}

	logger.Trace().Msg("building redis command")
	command := e.redis.B().JsonGet().Key(string(entityId)).Paths(componentPath(componentName)).Build()

	logger.Trace().Msg("executing redis command")
	result := e.redis.Do(context.Background(), command)

	logger.Trace().Msg("checking redis result")
	if result.Error() != nil {
		return result.Error()
	}

	err = result.DecodeJSON(component)

	if err != nil {
		return err
	}

	return nil
}

// HasEntity returns true if the entity exists in the engine. Any id passed will be
func (e *Engine) HasEntity(id uid.UID) (bool, error) {
	exists, err := e.redis.Do(context.Background(), e.redis.B().Exists().Key(string(id)).Build()).AsBool()

	if err != nil {
		return false, err
	}

	return exists, nil
}

// HasComponent returns true if the component exists on the entity.
func (e *Engine) HasComponent(entityId uid.UID, componentName string) (bool, error) {
	res, err := e.redis.Do(
		context.Background(),
		e.redis.B().JsonGet().Key(string(entityId)).Paths(componentPath(componentName)).Build(),
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

// RemoveComponent removes the named component from the given entity. This will trigger the
// `events.ComponentRemovedEvent` event. If the entity or component does not exist, an error will be returned. This
// method requires that an empty component be passed in, this will be used to unmarshal and add the current component
// value to the event.
func (e *Engine) RemoveComponent(entityId uid.UID, componentName string, valueType interface{}) error {
	logger := e.logger.
		With().
		Str("component", "entities").
		Str("entityId", string(entityId)).
		Str("componentName", componentName).
		Logger()
	logger.Debug().Msg("removing component")

	logger.Trace().Msg("checking if entity exists")
	exists, err := e.HasEntity(entityId)

	if err != nil {
		return err
	}

	if !exists {
		logger.Error().Msg("entity does not exist")
		return engineErrors.EntityNotFoundError{
			Id: entityId,
		}
	}

	logger.Trace().Msg("checking if component exists")
	exists, err = e.HasComponent(entityId, componentName)

	if err != nil {
		return err
	}

	if !exists {
		logger.Error().Msg("component does not exist")
		return engineErrors.ComponentNotFoundError{
			EntityId: string(entityId),
			Name:     componentName,
		}
	}

	logger.Trace().Msg("getting previous component value")
	err = e.GetComponent(entityId, componentName, &valueType)

	if err != nil {
		return err
	}

	logger.Trace().Msg("building redis commands")
	commands := rueidis.Commands{
		e.redis.B().JsonDel().Key(string(entityId)).Path(componentPath(componentName)).Build(),
	}

	logger.Trace().Msg("building publish commands for events")
	commands = append(commands, e.GetPublishCommandsForEvents(
		buildComponentRemovedEvent(entityId, componentName, valueType))...)

	logger.Trace().Msg("executing redis commands")
	results := e.redis.DoMulti(
		context.Background(),
		commands...,
	)

	logger.Trace().Msg("checking redis results")
	for _, result := range results {
		if result.Error() != nil {
			logger.Error().Err(result.Error()).Msg("error removing component")
		}
	}

	logger.Trace().Msg("component removed")

	return nil
}

// UpdateComponent updates a component on an entity. This will trigger the `events.ComponentUpdatedEvent` event to be
// published. If the entity does not exist, an error will be returned. If the component does not exist, an error will be
// returned.
func (e *Engine) UpdateComponent(entityId uid.UID, componentName string, component interface{}) error {
	logger := e.logger.With().
		Str("component", "entities").
		Str("entityId", string(entityId)).
		Str("componentName", componentName).Logger()

	logger.Debug().Msg("updating component")

	logger.Trace().Msg("checking if entity exists")
	exists, err := e.HasEntity(entityId)

	if err != nil {
		return err
	}

	if !exists {
		logger.Error().Msg("entity does not exist")
		return engineErrors.EntityNotFoundError{Id: entityId}
	}

	logger.Trace().Msg("checking if component exists")
	exists, err = e.HasComponent(entityId, componentName)

	if err != nil {
		return err
	}

	if !exists {
		logger.Error().Msg("component does not exist")
		return engineErrors.ComponentNotFoundError{EntityId: string(entityId), Name: componentName}
	}

	logger.Trace().Msg("getting previous value")
	prev := component

	err = e.GetComponent(entityId, componentName, &prev)

	if err != nil {
		return err
	}

	logger.Trace().Msg("building redis commands")
	commands := rueidis.Commands{
		e.redis.
			B().
			JsonSet().
			Key(string(entityId)).
			Path(componentPath(componentName)).
			Value(rueidis.JSON(component)).
			Build(),
	}

	logger.Trace().Msg("building publish commands for events")
	commands = append(commands, e.GetPublishCommandsForEvents(
		engineEvents.ComponentUpdatedEvent{
			EntityId:      entityId,
			Name:          componentName,
			Value:         component,
			PreviousValue: prev,
		},
	)...)

	logger.Trace().Msg("executing redis commands")

	results := e.redis.DoMulti(
		context.Background(),
		commands...,
	)

	cue := engineErrors.UpdateComponentErrors{}

	logger.Trace().Msg("checking redis results")
	for _, result := range results {
		if result.Error() != nil {
			cue.Add(result.Error())
		}
	}

	if cue.HasErrors() {
		return cue
	}

	return nil
}

func buildComponentAddedEvents(entityId uid.UID, components map[string]interface{}) []Event {
	events := make([]Event, 0)

	for name, value := range components {
		events = append(events, engineEvents.ComponentAddedEvent{
			EntityId: entityId,
			Name:     name,
			Value:    value,
		})
	}

	return events
}

func buildComponentRemovedEvent(entityId uid.UID, componentName string, value interface{}) Event {
	return engineEvents.ComponentRemovedEvent{
		EntityId: entityId,
		Name:     componentName,
		Value:    value,
	}
}

func buildEntityAddedEvent(entityId uid.UID) Event {
	return engineEvents.EntityAddedEvent{
		Id: entityId,
	}
}

func buildEntityAndComponentAddedEvents(entityId uid.UID, components map[string]interface{}) []Event {
	events := []Event{
		buildEntityAddedEvent(entityId),
	}
	return append(events, buildComponentAddedEvents(entityId, components)...)
}

func getComponentMap(entity interface{}) map[string]interface{} {
	return structs.Map(entity)
}

func componentPath(componentName string) string {
	return fmt.Sprintf(".Entity.%s", componentName)
}

func entityIsStruct(entity interface{}) bool {
	if entity == nil {
		return false
	}

	return structs.IsStruct(entity)
}
