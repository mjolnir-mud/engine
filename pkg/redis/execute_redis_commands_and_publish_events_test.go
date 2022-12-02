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

package redis

import (
	"context"
	"testing"

	"github.com/mjolnir-engine/engine/internal/redis"
	internalTesting "github.com/mjolnir-engine/engine/internal/testing"
	"github.com/mjolnir-engine/engine/pkg/event"
	"github.com/rueian/rueidis"
	"github.com/stretchr/testify/assert"
)

func TestExecuteRedisCommandsAndPublishEvents(t *testing.T) {
	ctx := internalTesting.Setup(t)
	defer internalTesting.Teardown(t, ctx)

	ctx = redis.WithRedis(ctx)
	client := GetClient(ctx)

	received := make(chan *fakeEvent)

	ctx, id := Subscribe(ctx, fakeEvent{}, func(event event.Message) {
		ev := &fakeEvent{}

		err := event.Unmarshal(ev)

		assert.NoError(t, err)

		received <- ev
	})
	defer Unsubscribe(ctx, id)

	ExecuteRedisCommandsAndPublishEvents(ctx, rueidis.Commands{
		client.B().Set().Key("test").Value("test").Build(),
	}, []event.Event{fakeEvent{}})

	assert.Equal(t, &fakeEvent{}, <-received)

	result := client.Do(ctx, client.B().Get().Key("test").Build())

	assert.NoError(t, result.Error())
	s, _ := result.ToString()
	assert.Equal(t, "test", s)
}

func TestMustExecuteRedisCommandsAndPublishEvents(t *testing.T) {
	ctx := internalTesting.Setup(t)
	defer internalTesting.Teardown(t, ctx)

	ctx = redis.WithRedis(ctx)
	client := GetClient(ctx)

	received := make(chan *fakeEvent)

	ctx, id := Subscribe(ctx, fakeEvent{}, func(event event.Message) {
		ev := &fakeEvent{}

		err := event.Unmarshal(ev)

		assert.NoError(t, err)

		received <- ev
	})
	defer Unsubscribe(ctx, id)

	MustExecuteRedisCommandsAndPublishEvents(ctx, rueidis.Commands{
		client.B().Set().Key("test").Value("test").Build(),
	}, []event.Event{fakeEvent{}})

	assert.Equal(t, &fakeEvent{}, <-received)

	result := client.Do(ctx, client.B().Get().Key("test").Build())

	assert.NoError(t, result.Error())
	s, _ := result.ToString()
	assert.Equal(t, "test", s)

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	assert.Panics(t, func() {
		MustExecuteRedisCommandsAndPublishEvents(ctx, rueidis.Commands{
			client.B().Set().Key("test").Value("test").Build(),
		}, []event.Event{fakeEvent{}})
	})
}
