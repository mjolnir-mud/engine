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
)

// UpdateComponentErrors is a collection of errors that occurred while adding components to an entity.
type UpdateComponentErrors struct {
	Errors []string
}

func (e UpdateComponentErrors) Error() string {
	errorStrings := make([]string, len(e.Errors))

	for i, err := range e.Errors {
		errorStrings[i] = err
	}

	return fmt.Sprintf("%d error(s) occurred while updating the component", len(e.Errors))
}

func (e *UpdateComponentErrors) Add(err error) {
	e.Errors = append(e.Errors, err.Error())
}

func (e UpdateComponentErrors) HasErrors() bool {
	return len(e.Errors) > 0
}
