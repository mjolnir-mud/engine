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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrContextValue(t *testing.T) {
	ctx := context.Background()
	_ = os.Setenv("bar", "baz")
	defer func() {
		_ = os.Unsetenv("bar")
	}()

	ctx = WithEnvironmentVariableOrDefault(ctx, "foo", "bar")

	assert.Equal(t, "bar", GetEnvOrContextValue(ctx, "foo"))
	assert.Equal(t, "baz", GetEnvOrContextValue(ctx, "bar"))
}

func TestMustGetEnvOrContextValue(t *testing.T) {
	ctx := context.Background()
	_ = os.Setenv("bar", "baz")
	defer func() {
		_ = os.Unsetenv("bar")
	}()

	ctx = WithEnvironmentVariableOrDefault(ctx, "foo", "bar")

	assert.Equal(t, "bar", MustGetEnvOrContextValue(ctx, "foo"))
	assert.Equal(t, "baz", MustGetEnvOrContextValue(ctx, "bar"))

	assert.Panics(t, func() {
		MustGetEnvOrContextValue(ctx, "baz")
	})
}

func TestGetEnvOrContextValueAsInt(t *testing.T) {
	ctx := context.Background()
	_ = os.Setenv("bar", "1")
	defer func() {
		_ = os.Unsetenv("bar")
	}()

	ctx = WithEnvironmentVariableOrDefault(ctx, "foo", "2")

	assert.Equal(t, 2, GetEnvOrContextValueAsInt(ctx, "foo"))
	assert.Equal(t, 1, GetEnvOrContextValueAsInt(ctx, "bar"))
}

func TestMustGetEnvOrContextValueAsInt(t *testing.T) {
	ctx := context.Background()
	_ = os.Setenv("bar", "1")
	defer func() {
		_ = os.Unsetenv("bar")
	}()

	ctx = WithEnvironmentVariableOrDefault(ctx, "foo", "2")

	assert.Equal(t, 2, MustGetEnvOrContextValueAsInt(ctx, "foo"))
	assert.Equal(t, 1, MustGetEnvOrContextValueAsInt(ctx, "bar"))

	assert.Panics(t, func() {
		MustGetEnvOrContextValueAsInt(ctx, "baz")
	})
}
