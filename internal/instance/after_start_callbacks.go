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

package instance

// AfterStart - called after the engine starts.
var afterStartCallbacks []func()

// AfterStartForEnv - called after the engine starts for a specific environment
var afterStartCallbacksForEnv map[string][]func()

// AfterServiceStart - called after a service starts
var afterServiceStartCallbacks map[string][]func()

// AfterServiceStartForEnv - called after a service starts for a specific environment
var afterServiceStartCallbacksForEnv map[string]map[string][]func()

func initializeAfterStartCallbacks() {
	afterStartCallbacks = make([]func(), 0)
	afterStartCallbacksForEnv = make(map[string][]func())
	afterServiceStartCallbacks = make(map[string][]func())
	afterServiceStartCallbacksForEnv = make(map[string]map[string][]func())
}

func RegisterAfterStartCallback(callback func()) {
	afterStartCallbacks = append(afterStartCallbacks, callback)
}

func RegisterAfterStartCallbackForEnv(env string, callback func()) {
	afterStartCallbacksForEnv[env] = append(afterStartCallbacksForEnv[env], callback)
}

func RegisterAfterServiceStartCallback(service string, callback func()) {
	_, ok := afterServiceStartCallbacks[service]

	if !ok {
		afterServiceStartCallbacks[service] = make([]func(), 0)
	}

	afterServiceStartCallbacks[service] = append(afterServiceStartCallbacks[service], callback)
}

func RegisterAfterServiceStartCallbackForEnv(service string, env string, callback func()) {
	_, ok := afterServiceStartCallbacksForEnv[service]

	if !ok {
		afterServiceStartCallbacksForEnv[service] = make(map[string][]func())
	}

	_, ok = afterServiceStartCallbacksForEnv[service][env]

	if !ok {
		afterServiceStartCallbacksForEnv[service][env] = make([]func(), 0)
	}

	afterServiceStartCallbacksForEnv[service][env] = append(afterServiceStartCallbacksForEnv[service][env], callback)
}

func callAfterStartCallbacks() {
	log.Info().Msg("Calling after start callbacks")
	for _, callback := range afterStartCallbacks {
		callback()
	}
}

func callAfterStartCallbacksForEnv(env string) {
	log.Info().Msg("Calling after start callbacks for env " + env)
	for _, callback := range afterStartCallbacksForEnv[env] {
		callback()
	}
}

func callAfterServiceStartCallbacks(service string) {
	log.Info().Msg("Calling after service start callbacks for " + service)
	for _, callback := range afterServiceStartCallbacks[service] {
		callback()
	}
}

func callAfterServiceStartCallbacksForEnv(service string, env string) {
	log.Info().Msg("Calling after service start callbacks for " + service + " and env " + env)
	for _, callback := range afterServiceStartCallbacksForEnv[service][env] {
		callback()
	}
}
