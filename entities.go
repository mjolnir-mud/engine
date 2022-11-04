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
	"reflect"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
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

// AddEntity adds an entity to the engine. If the entity id is not set, it will be set to a new UUID. If the entity
// already exists, an error will be returned.
func (e *Engine) AddEntity(entity interface{}) error {
	err := typeCheckStructPointer(entity)

	if err != nil {
		return err
	}

	logger := e.logger.With().Str("component", "entities").Logger()
	logger.Debug().Msg("adding entity")

	logger.Trace().Msg("prepping id")
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

	logger = logger.With().Str("entityId", string(id)).Logger()
	logger.Trace().Msg("checking if entity exists")
	exists, err := e.HasEntity(entity)

	if err != nil {
		return err
	}

	if exists {
		return engineErrors.EntityExistsError{
			Id: id,
		}
	}

	logger.Trace().Msg("building redis command")
	record := EntityRecord{
		Id:      string(id),
		Version: 1,
		Entity:  entity,
	}

	commands := rueidis.Commands{
		e.redis.B().JsonSet().Key(string(id)).Path(".").Value(rueidis.JSON(record)).Build(),
	}

	components := getComponentMap(entity)
	commands = append(commands, e.GetPublishCommandsForEvents(
		buildEntityAndComponentAddedEvents(id, components)...)...,
	)

	results := e.redis.DoMulti(
		context.Background(),
		commands...,
	)

	for _, result := range results {
		if result.Error() != nil {
			logger.Error().Err(result.Error()).Msg("error adding entity")
		}
	}

	logger.Trace().Msg("entity added")

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

// GetEntity populates an entity with data from the engine. It requires a struct to be passed that has at least an `Id`
// field. If the entity does not exist, an error will be returned. If the struct does not have an `Id` field, an error
// will be returned.
func (e *Engine) GetEntity(entity interface{}) error {
	logger := e.logger.With().Str("component", "entities").Logger()

	logger.Debug().Msg("getting entity")
	logger.Trace().Msg("getting entity id")
	id, err := getEntityId(entity)

	if err != nil {
		return err
	}

	logger = logger.With().Str("entityId", string(id)).Logger()
	logger.Trace().Msg("checking if entity exists")

	exists, err := e.HasEntity(entity)

	if err != nil {
		return err
	}

	if !exists {
		logger.Error().Msg("entity does not exist")
		return engineErrors.EntityNotFoundError{
			Id: id,
		}
	}

	logger.Trace().Msg("building redis command")
	command := e.redis.B().JsonGet().Key(string(id)).Paths(".").Build()

	logger.Trace().Msg("executing redis command")
	result := e.redis.Do(context.Background(), command)

	logger.Trace().Msg("checking redis result")
	if result.Error() != nil {
		return result.Error()
	}

	logger.Trace().Msg("decoding redis result")
	err = result.DecodeJSON(entity)

	if err != nil {
		return err
	}

	return nil
}

// HasEntity returns true if the entity exists in the engine. It accepts an entity struct or an id.
func (e *Engine) HasEntity(entityOrId interface{}) (bool, error) {
	id, err := getEntityId(entityOrId)

	if err != nil {
		return false, err
	}

	exists, err := e.redis.Do(context.Background(), e.redis.B().Exists().Key(string(id)).Build()).AsBool()

	if err != nil {
		return false, err
	}

	return exists, nil
}

// HasComponent returns true if the component exists on the entity.
func (e *Engine) HasComponent(entityOrId interface{}, componentName string) (bool, error) {
	entityId, err := getEntityId(entityOrId)

	if err != nil {
		return false, err
	}

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

// RemoveEntity removes the entity from the engine. If the entity does not exist, an error will be returned. Each
// component will be removed from the engine as well.
func (e *Engine) RemoveEntity(id uid.UID) error {
	//logger := e.logger.With().Str("component", "entities").Str("entityId", string(id)).Logger()
	//
	//logger.Debug().Msg("removing entity")
	//
	//logger.Trace().Msg("checking if entity exists")
	//entity, err := e.GetEntity(id)
	//
	//logger.Trace().Msg("building redis commands")
	//commands := rueidis.Commands{
	//	e.redis.B().Del().Key(string(id)).Build(),
	//}
	//
	//components := getComponentMap()

	return nil
}

// UpdateEntity updates a component on an entity. This will trigger the `events.ComponentUpdatedEvent` event to be
// published. If the entity does not exist, an error will be returned. If the component does not exist, an error will be
// returned.
func (e *Engine) UpdateEntity(entity interface{}) error {
	err := typeCheckStructPointer(entity)

	if err != nil {
		return err
	}

	m := getComponentMap(entity)

	var id uid.UID
	id, ok := m["Id"].(uid.UID)

	if !ok {
		return engineErrors.EntityNotFoundError{
			Id: m["Id"].(uid.UID),
		}
	}

	logger := e.logger.With().
		Str("component", "entities").
		Str("entityId", string(id)).
		Logger()

	logger.Debug().Msg("updating component")

	logger.Trace().Msg("checking if entity exists")

	entityRecord, err := e.getEntityRecord(id)

	if err != nil {
		return err
	}

	oldEntityMap := entityRecord.Entity.(map[string]interface{})
	newEntityMap := getComponentMap(entity)

	if err != nil {
		return err
	}

	logger.Trace().Msg("comparing old and new entity")

	if reflect.DeepEqual(oldEntityMap, newEntityMap) {
		logger.Trace().Msg("entity has not changed")
		return nil
	}

	entityRecord.Version++
	entityRecord.Entity = entity

	// set the updated entity
	commands := rueidis.Commands{
		e.redis.B().JsonSet().Key(string(id)).Path(".").Value(rueidis.JSON(entityRecord)).Build(),
	}

	logger.Trace().Msg("building publish commands for events")

	commands = append(commands, e.GetPublishCommandsForEvents(
		buildComponentUpdatedEvents(newEntityMap, oldEntityMap)...,
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

func typeCheckStructPointer(entity interface{}) error {
	if entity == nil {
		return engineErrors.EntityDataTypedError{
			Expected: "not nil",
		}
	}

	entityType := reflect.TypeOf(entity)

	if entityType.Kind() != reflect.Ptr {
		return engineErrors.EntityDataTypedError{
			Expected: "pointer",
		}
	}

	if entityType.Elem().Kind() != reflect.Struct {
		return engineErrors.EntityDataTypedError{
			Expected: "struct",
		}
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

func buildComponentUpdatedEvents(entity map[string]interface{}, oldEntity map[string]interface{}) []Event {
	id, ok := entity["Id"].(uid.UID)

	if !ok {
		panic("unable to get entity id")
	}

	events := make([]Event, 0)

	for name, value := range entity {
		if oldEntity[name] == nil {
			events = append(events, engineEvents.ComponentAddedEvent{
				EntityId: id,
				Name:     name,
				Value:    value,
			})
		}

		if !reflect.DeepEqual(value, oldEntity[name]) {
			events = append(events, engineEvents.ComponentUpdatedEvent{
				EntityId:      id,
				Name:          name,
				Value:         value,
				PreviousValue: oldEntity[name],
			})
		}
	}

	for name, value := range oldEntity {
		if entity[name] == nil {
			events = append(events, engineEvents.ComponentRemovedEvent{
				EntityId: id,
				Name:     name,
				Value:    value,
			})
		}
	}

	return events
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

func getEntityId(entityOrId interface{}) (uid.UID, error) {
	if entityOrId == nil {
		return "", engineErrors.EntityDataTypedError{}
	}

	switch entityOrId.(type) {
	case uid.UID:
		return entityOrId.(uid.UID), nil
	default:
		m := getComponentMap(entityOrId)

		if id, ok := m["Id"]; ok {
			return id.(uid.UID), nil
		}
	}

	return uid.New(), engineErrors.EntityInvalidError{}
}

func componentPath(componentName string) string {
	return fmt.Sprintf(".Entity.%s", componentName)
}

func (e *Engine) getEntityRecord(id uid.UID) (EntityRecord, error) {
	var record EntityRecord

	result := e.redis.Do(
		context.Background(),
		e.redis.B().JsonGet().Key(string(id)).Build(),
	)

	if result.Error() != nil {
		return record, result.Error()
	}
	result.DecodeJSON(&record)
	result.DecodeJSON(&record)

	return record, nil
}
