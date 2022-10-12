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

package errors

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddComponentErrors_Error(t *testing.T) {
	e := AddComponentErrors{
		Errors: []string{
			"error 1",
		},
	}

	assert.Equal(t, "1 error(s) occurred while adding the component", e.Error())
}

func TestAddComponentErrors_Errors(t *testing.T) {
	e := AddComponentErrors{
		Errors: []string{
			"error 1",
		},
	}

	assert.Equal(t, []string{"error 1"}, e.Errors)
}

func TestAddComponentErrors_HasErrors(t *testing.T) {
	e := AddComponentErrors{
		Errors: []string{
			"error 1",
		},
	}

	assert.True(t, e.HasErrors())
}

func TestAddComponentErrors_Add(t *testing.T) {
	e := AddComponentErrors{
		Errors: make([]string, 0),
	}

	e.Add(fmt.Errorf("error 1"))

	assert.Equal(t, "1 error(s) occurred while adding the component", e.Error())
}
