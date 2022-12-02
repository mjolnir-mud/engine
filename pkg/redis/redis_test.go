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
)

type fakeEvent struct {
	Value string
}

func (e fakeEvent) Topic() string {
	return "engine:fake:event"
}

func (e fakeEvent) AllTopics() string {
	return "engine:*:event"
}

func setup(t *testing.T) context.Context {
	ctx := internalTesting.Setup(t)
	ctx = redis.WithRedis(ctx)

	return ctx
}

func teardown(t *testing.T, ctx context.Context) {
	internalTesting.Teardown(t, ctx)
}