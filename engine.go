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

package engine

import (
	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/rs/zerolog"
	"github.com/rueian/rueidis"
)

// Engine is an instance of the Mjolnir game engine.
type Engine struct {
	redis      rueidis.Client
	instanceId string

	// Logger is the logger for the engine. Plugins can use this logger to create their own tagged loggers. See
	// [zerolog](https://github.com/rs/zerolog) for more information.
	Logger zerolog.Logger
}

// New creates a new instance of the Mjolnir game engine. If the redis connection fails, an error is returned.
func New(config *Configuration) (*Engine, error) {
	redisClient, err := redis.New(config.Redis)

	if err != nil {
		return nil, err
	}

	return &Engine{
		redis:      redisClient,
		instanceId: config.InstanceId,
		Logger:     newLogger(config.Log),
	}, nil
}
