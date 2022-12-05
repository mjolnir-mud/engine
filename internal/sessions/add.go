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

	"github.com/mjolnir-engine/engine/pkg/events"
	"github.com/mjolnir-engine/engine/pkg/logger"
	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/mjolnir-engine/engine/pkg/uid"
)

func Add(ctx context.Context, id uid.UID) context.Context {
	l := logger.Get(ctx).With().Str("session_id", string(id)).Logger()
	l.Debug().Msg("adding session")

	sessCtx := context.WithValue(ctx, "id", id)
	sessCtx, stop := context.WithCancel(ctx)
	sessCtx = context.WithValue(sessCtx, StopFuncKey, stop)

	l.Trace().Msg("setting stop subscription")
	sessCtx, stopSid := redis.Subscribe(
		sessCtx,
		events.SessionRemoveEvent{
			Id: id,
		},
		func(e events.Message) {
			ev := events.SessionRemoveEvent{}
			err := e.Unmarshal(ev)

			if err != nil {
				panic(err)
			}

			stop()
		},
	)

	l.Trace().Msg("setting up receive subscription")
	sessCtx, receiveSid := redis.Subscribe(
		sessCtx,
		events.SessionReceiveDataEvent{
			Id: id,
		},
		func(e events.Message) {
			ev := events.SessionReceiveDataEvent{}
			err := e.Unmarshal(ev)

			if err != nil {
				panic(err)
			}
		},
	)

	l.Trace().Msg("setting up unsubscribe go-routine")
	go func() {
		<-sessCtx.Done()

		redis.Unsubscribe(sessCtx, stopSid)
		redis.Unsubscribe(sessCtx, receiveSid)
	}()

	sessionContexts := ctx.Value(ContextsKey).(map[uid.UID]context.Context)
	sessionContexts[id] = sessCtx
	ctx = context.WithValue(ctx, ContextsKey, sessionContexts)

	l.Trace().Msg("publishing session added event")
	redis.Publish(ctx, events.SessionAddedEvent{
		Id: id,
	})

	l.Debug().Msg("session added")

	return ctx
}
