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

package _refactor

import (
	"fmt"
	"github.com/mjolnir-mud/engine"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testEvent struct {
	Value string
}

func (e testEvent) Topic() string {

	return fmt.Sprintf("testing.%s", e.Value)
}

func setup() {
	viper.Set("env", "testing")
	Start("localhost", 6379, 1)
}

func teardown() {
	Stop()
}

func TestSubscribe(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan interface{})

	s := NewSubscription(&testEvent{Value: "testing"}, func(payload engine.EventPayload) {
		p := &testEvent{}

		_ = payload.Unmarshal(p)

		ch <- p
	})

	err := Publish(&testEvent{
		Value: "testing",
	})

	assert.Nil(t, err)

	v := <-ch

	assert.Equal(t, &testEvent{
		Value: "testing",
	}, v)

	s.Stop()
}

func TestPSubscribe(t *testing.T) {
	setup()
	defer teardown()
	ch := make(chan interface{})

	s := NewPatternSubscription(&testEvent{Value: "*"}, func(payload engine.EventPayload) {
		e := &testEvent{}
		_ = payload.Unmarshal(e)

		ch <- e
	})

	err := Publish(&testEvent{
		Value: "testing",
	})

	assert.Nil(t, err)

	v := <-ch

	assert.Equal(t, &testEvent{
		Value: "testing",
	}, v)

	s.Stop()
}
