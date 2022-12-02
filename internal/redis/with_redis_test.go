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
	"testing"

	"github.com/mjolnir-engine/engine/internal/engine_context"
	internalTesting "github.com/mjolnir-engine/engine/internal/testing"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rueian/rueidis"
	"github.com/stretchr/testify/assert"
)

func TestWithRedis(t *testing.T) {
	ctx := internalTesting.Setup(t)
	defer func() {
		ctx.Value(ClientContextKey).(rueidis.Client).Close()
		internalTesting.Teardown(t, ctx)
	}()

	ctx = WithRedis(ctx)

	assert.NotEmpty(t, engine_context.MustGetEnvOrContextValue(ctx, HostEnvVar))
	assert.NotEmpty(t, engine_context.MustGetEnvOrContextValue(ctx, PortEnvVar))
	assert.NotEmpty(t, engine_context.MustGetEnvOrContextValue(ctx, DatabaseEnvVar))

	assert.NotEmpty(t, ctx.Value(ClientContextKey))
	assert.NotEmpty(t, ctx.Value(LoggerContextKey))
	assert.Equal(t, map[uid.UID]*Subscription{}, ctx.Value(SubscriptionsContextKey))
}
