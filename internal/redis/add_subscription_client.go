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
	"strings"

	"github.com/mjolnir-engine/engine/pkg/event"
	"github.com/mjolnir-engine/engine/pkg/uid"
	"github.com/rueian/rueidis"
)

func AddSubscription(ctx context.Context, topic string, callback func(event event.Message)) (context.Context, uid.UID) {
	id := uid.New()
	client := ctx.Value(ClientContextKey).(rueidis.Client)
	dedicated, cancel := client.Dedicate()
	logger := GetLogger(ctx)

	logger.Debug().Str("topic", topic).Str("id", string(id)).Msg("adding subscription")

	setClientPubSubHooks(dedicated, cancel, callback)

	if strings.Contains(topic, "*") {
		dedicated.Do(ctx, dedicated.B().Psubscribe().Pattern(topic).Build())
	} else {
		dedicated.Do(ctx, dedicated.B().Subscribe().Channel(topic).Build())
	}

	subs := ctx.Value(SubscriptionsContextKey).(map[uid.UID]*Subscription)

	subs[id] = &Subscription{
		client: dedicated,
		topic:  topic,
		cancel: cancel,
		ctx:    ctx,
	}

	ctx = context.WithValue(ctx, SubscriptionsContextKey, subs)

	return ctx, id
}

func setClientPubSubHooks(client rueidis.DedicatedClient, cancel func(), callback func(event.Message)) {
	go func() {
		defer cancel()

		wait := client.SetPubSubHooks(rueidis.PubSubHooks{
			OnMessage: func(m rueidis.PubSubMessage) {
				callback(event.Message(m))
			},
		})

		<-wait
	}()
}
