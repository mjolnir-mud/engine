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
	"github.com/rs/zerolog"
	"io"
	"os"
)

// LogConfiguration represents the Mjolnir logging configuration.
type LogConfiguration struct {
	// Level is the logging level.
	Level zerolog.Level
	Writer io.Writer
}

// newLogger creates a new logger. It accepts a `LogConfiguration` struct as its only argument.
func newLogger(config *LogConfiguration) zerolog.Logger {
	if config == nil {
		config = &LogConfiguration{
			Level: zerolog.TraceLevel,
			Writer: os.Stdout,
		}
	}

	if config.Writer == nil {
		config.Writer = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	if config.Level == zerolog.NoLevel {
		config.Level = zerolog.InfoLevel
	}

	return zerolog.New(config.Writer).With().Timestamp().Logger().Level(config.Level)
}
