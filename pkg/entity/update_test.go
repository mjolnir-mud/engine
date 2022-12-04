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

func TestUpdate(t *testing.T) {
	ctx := setup(t)
	defer teardown(t, ctx)

	t.Run("entity exists", func(t *testing.T) {
		fe := &fakeEntity{
			Id:    uid.New(),
			Value: "test",
		}

		err := Add(ctx, fe)
		assert.NoError(t, err)

		updated := make(chan *events.ComponentUpdatedEvent, 1)

		var id uid.UID
		ctx, id = redis.Subscribe(ctx, events.ComponentUpdatedEvent{
			EntityId: fe.Id,
			Name:     "Value",
		}, func(e event.Message) {
			ev := &events.ComponentUpdatedEvent{}

			err := e.Unmarshal(ev)
			assert.NoError(t, err)

			go func() { updated <- ev }()
		})
		defer redis.Unsubscribe(ctx, id)

		fe.Value = "updatedTest"

		err = Update(ctx, fe)

		u := <-updated

		assert.Equal(t, "updatedTest", u.Value)
		assert.Equal(t, fe.Id, u.EntityId)
		assert.Equal(t, "Value", u.Name)
		assert.Equal(t, "test", u.PreviousValue)
	})

	t.Run("entity does not exist", func(t *testing.T) {
		fe := &fakeEntity{
			Id:    uid.New(),
			Value: "test",
		}

		err := Update(ctx, fe)

		assert.Error(t, err)
	})
}
