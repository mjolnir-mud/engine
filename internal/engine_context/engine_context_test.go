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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ctx := New("test")
	assert.NotNil(t, ctx)

	assert.NotEmptyf(t, ctx.Value("instanceId"), "expected context id to be set")
	assert.Equal(t, ctx.Value("services"), []string{})
	assert.NotEmpty(t, ctx.Value("service"), "expected service to be set")
	assert.NotEmpty(t, ctx.Value("logger"), "expected logger to be set")

	// MONGO CONFIG
	assert.Equal(t, GetEnvOrContextValue(ctx, "MJOLNIR_ENV"), "development")
	assert.Equal(t, GetEnvOrContextValue(ctx, "MJOLNIR_LOG_LEVEL"), "debug")
	assert.Equal(t, GetEnvOrContextValue(ctx, "MJOLNIR_MONGO_HOST"), "localhost")
	assert.Equal(t, GetEnvOrContextValue(ctx, "MJOLNIR_MONGO_PORT"), "27017")
	assert.Equal(t, GetEnvOrContextValue(ctx, "MJOLNIR_MONGO_DATABASE"), "mjolnir_development")

	// REDIS CONFIG
}
