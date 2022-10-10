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
	"github.com/mjolnir-engine/engine/uid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeDataSource struct{}

func (f fakeDataSource) Name() string {
	return "fake"
}

func (f fakeDataSource) FindOne(map[string]interface{}) (interface{}, error) {
	panic("implement me")
}

func (f fakeDataSource) Find(map[string]interface{}) ([]interface{}, error) {
	panic("implement me")
}

func (f fakeDataSource) Save(interface{}) (*uid.UID, error) {
	panic("implement me")
}

func (f fakeDataSource) Count(map[string]interface{}) (int64, error) {
	panic("implement me")
}

func (f fakeDataSource) Delete(map[string]interface{}) error {
	panic("implement me")
}

func (f fakeDataSource) All() ([]interface{}, error) {
	panic("implement me")
}

func (f fakeDataSource) Start() error {
	panic("implement me")
}

func (f fakeDataSource) Stop() error {
	panic("implement me")
}

func TestEngine_RegisterDataSource(t *testing.T) {
	e := createEngineInstance()
	e.RegisterDataSource(fakeDataSource{})

	assert.Len(t, e.dataSourceRegistry.dataSources, 1)
}
