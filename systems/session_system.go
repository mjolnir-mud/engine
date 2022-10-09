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

package systems

type sessionSystem struct{}

func (e sessionSystem) Name() string {
	return "session"
}

func (e sessionSystem) Component() string {
	return "session"
}

func (e sessionSystem) Match(_ string, _ interface{}) bool {
	return true
}

func (e sessionSystem) ComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e sessionSystem) ComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (e sessionSystem) ComponentRemoved(_ string, _ string) error { return nil }

func (e sessionSystem) MatchingComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e sessionSystem) MatchingComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (e sessionSystem) MatchingComponentRemoved(_ string, _ string) error { return nil }

var Session = sessionSystem{}
