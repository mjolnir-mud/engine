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

package entity

import (
	"context"

	"github.com/mjolnir-engine/engine/pkg/errors"
	"github.com/mjolnir-engine/engine/pkg/logger"
	"github.com/mjolnir-engine/engine/pkg/redis"
	"github.com/rueian/rueidis"
)

// Get populates an entity with data from the engine. It requires a struct to be passed that has at least an `Id`
// field. If the entity does not exist, an error will be returned. If the struct does not have an `Id` field, an error
// will be returned.
func Get(ctx context.Context, entity interface{}) error {
	l := logger.Get(ctx).With().Str("component", "entities").Logger()

	l.Debug().Msg("getting entity")
	l.Trace().Msg("getting entity id")

	id := mustGetEntityId(entity)

	l = l.With().Str("entityId", string(id)).Logger()
	l.Trace().Msg("checking if entity exists")

	exists := HasEntity(ctx, entity)

	if !exists {
		l.Error().Msg("entity does not exist")
		return errors.EntityNotFoundError{
			Id: id,
		}
	}

	l.Trace().Msg("building redis command")
	redis.UnmarshallResult(redis.MustExecuteCommands(ctx, rueidis.Commands{
		redis.GetClient(ctx).B().JsonGet().Key(string(id)).Paths(".Entity").Build(),
	})[0], entity)

	l.Trace().Msg("decoding redis result")

	return nil
}
