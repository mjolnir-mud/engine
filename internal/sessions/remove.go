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

func Remove(ctx context.Context, id uid.UID) context.Context {
	l := logger.Get(ctx).With().Str("session_id", string(id)).Logger()

	l.Debug().Msg("removing session")
	sessionsContexts := ctx.Value(ContextsKey).(map[uid.UID]context.Context)
	sessCtx := sessionsContexts[id]

	if sessCtx.Err() != context.Canceled {
		l.Trace().Msg("cancelling session context, as it has not yet been cancelled")
		stop := sessCtx.Value(StopFuncKey).(context.CancelFunc)
		stop()
	}

	delete(sessionsContexts, id)

	l.Trace().Msg("publishing session removed event")
	redis.Publish(ctx, events.SessionRemovedEvent{Id: id})

	l.Debug().Msg("session removed")

	return ctx
}
