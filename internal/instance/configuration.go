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

import (
	"github.com/mjolnir-mud/engine"
)

var Configs map[string]func(configuration *engine.Configuration) *engine.Configuration

func initializeConfigurations() {
	Configs = map[string]func(configuration *engine.Configuration) *engine.Configuration{
		"default": func(*engine.Configuration) *engine.Configuration {
			return &engine.Configuration{
				Redis: engine.RedisConfiguration{
					Host: "localhost",
					Port: 6379,
					Db:   0,
				},
			}
		},
		"test": func(*engine.Configuration) *engine.Configuration {
			return &engine.Configuration{
				Redis: engine.RedisConfiguration{
					Host: "localhost",
					Port: 6379,
					Db:   1,
				},
			}
		},
	}
}

func callConfigureForEnv(env string) *engine.Configuration {
	defaultConfig := Configs["default"](&engine.Configuration{})
	configForEnv, ok := Configs[env]

	if !ok {
		panic("No configuration found for environment: " + env)
	}

	return configForEnv(defaultConfig)
}

func ConfigureForEnv(env string, cb func(config *engine.Configuration) *engine.Configuration) {
	Configs[env] = cb
}
