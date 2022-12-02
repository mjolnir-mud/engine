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
	"testing"

	"github.com/mjolnir-engine/engine/pkg/event"
	"github.com/mjolnir-engine/engine/pkg/events"
	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/stretchr/testify/assert"
)

func TestEngine_Add(t *testing.T) {
	ctx := setup(t)
	defer teardown(t, ctx)

	t.Run("it adds the entity and publishes the correct events", func(t *testing.T) {

		fe := &fakeEntity{
			Id:    uid.New(),
			Value: "test",
		}

		receivedEntityAdded := make(chan *events.EntityAddedEvent, 1)

		ctx, saddId := redis.PSubscribe(ctx, events.EntityAddedEvent{}, func(e event.Message) {
			ev := &events.EntityAddedEvent{}

			err := e.Unmarshal(ev)

			assert.NoError(t, err)

			go func() { receivedEntityAdded <- ev }()
		})
		defer redis.Unsubscribe(ctx, saddId)

		receivedCompnentAdded := make(chan *events.ComponentAddedEvent, 1)

		ctx, caddId := redis.PSubscribe(ctx, events.ComponentAddedEvent{}, func(e event.Message) {
			ev := &events.ComponentAddedEvent{}

			err := e.Unmarshal(ev)

			assert.NoError(t, err)

			go func() { receivedCompnentAdded <- ev }()
		})
		defer redis.Unsubscribe(ctx, caddId)

		err := Add(ctx, fe)

		assert.Nil(t, err)
		assert.Equal(t, &events.EntityAddedEvent{
			Id: fe.Id,
		}, <-receivedEntityAdded)

		componentEvents := make([]*events.ComponentAddedEvent, 0)
		for i := 0; i < 2; i++ {
			componentEvents = append(componentEvents, <-receivedCompnentAdded)
		}

		assert.Contains(t, componentEvents, &events.ComponentAddedEvent{
			EntityId: fe.Id,
			Name:     "Value",
			Value:    "test",
		})

		assert.Contains(t, componentEvents, &events.ComponentAddedEvent{
			EntityId: fe.Id,
			Name:     "Id",
			Value:    string(fe.Id),
		})
	})
}
