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
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakePlugin struct {
	StartCalled chan *Engine
	StopCalled  chan *Engine
}

func (f fakePlugin) Name() string {
	return "fake"
}

func (f fakePlugin) Start(e *Engine) error {
	go func() { f.StartCalled <- e }()

	return nil
}

func (f fakePlugin) Stop(e *Engine) error {
	go func() { f.StopCalled <- e }()

	return nil
}

func TestEngine_StartsPlugins(t *testing.T) {
	engine := createEngineInstance()
	fp := fakePlugin{
		StartCalled: make(chan *Engine, 1),
	}

	engine.RegisterPlugin(fp)

	engine.Start("test")
	defer engine.Stop()

	e := <-fp.StartCalled

	assert.Equal(t, engine, e)
}

func TestEngine_StopsPlugins(t *testing.T) {
	engine := createEngineInstance()
	fp := fakePlugin{
		StopCalled: make(chan *Engine, 1),
	}

	engine.RegisterPlugin(fp)

	engine.Start("test")
	engine.Stop()

	e := <-fp.StopCalled

	assert.Equal(t, engine, e)
}
