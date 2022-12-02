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
	"fmt"

	"github.com/mjolnir-engine/engine/internal/engine_context"
	"github.com/mjolnir-engine/engine/pkg/logger"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rueian/rueidis"
)

func WithRedis(ctx context.Context) context.Context {
	ctx = engine_context.WithEnvironmentVariableOrDefault(ctx, HostEnvVar, "localhost")
	ctx = engine_context.WithEnvironmentVariableOrDefault(ctx, PortEnvVar, "6379")
	ctx = engine_context.WithEnvironmentVariableOrDefault(ctx, DatabaseEnvVar, "0")

	redisHost := engine_context.MustGetEnvOrContextValue(ctx, HostEnvVar)
	redisPort := engine_context.MustGetEnvOrContextValueAsInt(ctx, PortEnvVar)
	redisDb := engine_context.MustGetEnvOrContextValueAsInt(ctx, DatabaseEnvVar)

	l := logger.GetLogger(ctx).
		With().
		Str("redisHost", redisHost).
		Int("redisPort", redisPort).
		Int("db", redisDb).
		Logger()

	l.Info().Msg("connecting to redis")

	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{
			fmt.Sprintf("%s:%d", redisHost, redisPort),
		},
		SelectDB: redisDb,
	})

	if err != nil {
		l.Panic().Err(err).Msg("failed to connect to redis")
	}

	ctx = context.WithValue(ctx, ClientContextKey, c)
	ctx = context.WithValue(ctx, LoggerContextKey, l)
	ctx = context.WithValue(ctx, SubscriptionsContextKey, make(map[uid.UID]*Subscription))

	return ctx
}
