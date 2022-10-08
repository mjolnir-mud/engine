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

package testing

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/internal/instance"
	"github.com/mjolnir-mud/engine/internal/redis"
)

var engineSetupCallbacks = make(map[string]func())
var engineSetupNames = make([]string, 0)
var engineTestRunning = false
var engineMux = make(chan bool)

func RegisterSetupCallback(plugin string, cb func()) {
	engineSetupNames = append(engineSetupNames, plugin)
	engineSetupCallbacks[plugin] = cb
}

func Setup(service string) chan bool {
	if !engineTestRunning {
		engineTestRunning = true
	} else {
		<-engineMux
		engineTestRunning = true
	}

	engine.Initialize("testing", "testing")

	engine.ConfigureForEnv("testing", func(cfg *engine.Configuration) *engine.Configuration {
		return cfg
	})

	ch := make(chan bool)
	engine.RegisterAfterStartCallback(func() {
		go func() { ch <- true }()
	})

	callBeforeStartCallbacks()

	engine.Start()

	instance.StartService(service)

	return ch
}

func callBeforeStartCallbacks() {
	for _, name := range engineSetupNames {
		engineSetupCallbacks[name]()
	}
}

func Teardown() {
	engineSetupCallbacks = make(map[string]func())
	flushed := make(chan error)
	go func() { flushed <- redis.FlushAll() }()

	<-flushed

	engine.Stop()
	engineTestRunning = false
	go func() { engineMux <- true }()
}
