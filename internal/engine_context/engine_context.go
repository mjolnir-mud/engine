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

package engine_context

import (
	"context"
	"fmt"
	"os"

	"github.com/fugufish/uncontext"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rs/zerolog"
)

var logLevels = map[string]zerolog.Level{
	"trace": zerolog.TraceLevel,
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
	"panic": zerolog.PanicLevel,
}

func New(service string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "service", service)
	ctx = WithEnvironmentVariableOrDefault(ctx, "MJOLNIR_ENV", "development")
	ctx = WithEnvironmentVariableOrDefault(ctx, "MJOLNIR_LOG_LEVEL", "debug")

	env := MustGetEnvOrContextValue(ctx, "MJOLNIR_ENV")

	// MONGO CONFIG
	ctx = WithEnvironmentVariableOrDefault(ctx, "MJOLNIR_MONGO_HOST", "localhost")
	ctx = WithEnvironmentVariableOrDefault(ctx, "MJOLNIR_MONGO_PORT", "27017")
	ctx = WithEnvironmentVariableOrDefault(ctx, "MJOLNIR_MONGO_DATABASE", fmt.Sprintf("mjolnir_%s", env))

	// REDIS CONFIG

	ctx = uid.WithUID(ctx, "instanceId", uid.New())
	ctx = context.WithValue(ctx, "services", make([]string, 0))
	ctx = context.WithValue(ctx, "service", service)

	ctx = withLogger(ctx)

	return ctx
}

func WithEnvironmentVariableOrDefault(ctx context.Context, key string, defaultValue string) context.Context {
	value := os.Getenv(key)

	if value == "" {
		value = defaultValue
	}

	return context.WithValue(ctx, key, value)
}

func withLogger(ctx context.Context) context.Context {
	logLevel := logLevels[MustGetEnvOrContextValue(ctx, "MJOLNIR_LOG_LEVEL")]
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(writer).With().
		Timestamp().
		Str("service", uncontext.MustGetString(ctx, "service")).
		Str("instanceId", string(uid.MustGetUID(ctx, "instanceId"))).
		Logger().Level(logLevel)

	return context.WithValue(ctx, "logger", logger)
}
