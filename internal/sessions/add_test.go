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

package sessions

import (
	"context"
	"testing"

	"github.com/mjolnir-engine/engine/pkg/events"
	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	ctx := setup(t)
	defer teardown(t, ctx)

	ctx = WithSessionsContext(ctx)
	id := uid.New()

	addedChan := make(chan events.SessionAddedEvent)

	ctx, sid := redis.Subscribe(ctx, events.SessionAddedEvent{
		Id: id,
	}, func(e events.Message) {
		ev := &events.SessionAddedEvent{}

		err := e.Unmarshal(ev)

		assert.NoError(t, err)

		go func() { addedChan <- *ev }()
	})
	defer redis.Unsubscribe(ctx, sid)

	ctx = Add(ctx, id)
	defer Remove(ctx, id)

	ev := <-addedChan
	assert.Equal(t, id, ev.Id)

	sessions := ctx.Value(ContextsKey).(map[uid.UID]context.Context)
	sessCtx := sessions[id]

	assert.NotNil(t, sessCtx)
}
