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

	"github.com/mjolnir-engine/engine/pkg/events"
	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/stretchr/testify/assert"
)

func TestRemoveEntity(t *testing.T) {
	ctx := setup(t)
	defer teardown(t, ctx)

	fe := &fakeEntity{
		Id:    uid.New(),
		Value: "test",
	}

	err := Add(ctx, fe)
	assert.NoError(t, err)

	removed := make(chan *events.EntityRemovedEvent, 1)

	var suid uid.UID
	ctx, suid = redis.Subscribe(
		ctx,
		events.EntityRemovedEvent{
			Id: fe.Id,
		},
		func(e events.Message) {
			ev := &events.EntityRemovedEvent{}

			err := e.Unmarshal(ev)

			assert.NoError(t, err)

			go func() { removed <- ev }()
		},
	)

	defer redis.Unsubscribe(ctx, suid)

	err = Remove(ctx, fe)
	ev := <-removed

	assert.NoError(t, err)
	assert.Equal(t, fe.Id, ev.Id)
}
