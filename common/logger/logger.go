/*
 * pd2mm
 * Copyright (C) 2025 pd2mm contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type MultiLogger struct {
	writers []io.Writer
	loggers []*log.Logger
}

var SharedLogger = NewMultiLogger(os.Stdout) //nolint:gochecknoglobals // allowed

// Create a new MultiLogger.
func NewMultiLogger(wrs ...io.Writer) *MultiLogger {
	loggers := new(MultiLogger)

	loggers.writers = make([]io.Writer, len(wrs))
	loggers.loggers = make([]*log.Logger, len(wrs))

	for i, w := range wrs {
		loggers.writers[i] = w
		loggers.loggers[i] = log.NewWithOptions(w, log.Options{ //nolint:exhaustruct // allowed
			ReportCaller:    false,
			ReportTimestamp: true,
			TimeFormat:      time.Kitchen,
		})
	}

	return loggers
}

// Register writers to the SharedLogger.
func RegisterLogger(wrs ...io.Writer) {
	SharedLogger = NewMultiLogger(wrs...)
}

// Write a DEBUG message with key value pairs.
func (ml *MultiLogger) Debug(msg any, kvs ...any) {
	for _, l := range ml.loggers {
		l.Debug(msg, kvs...)
	}
}

// Write a DEBUG message with fmt.Sprintf.
func (ml *MultiLogger) Debugf(format string, a ...any) {
	for _, l := range ml.loggers {
		l.Debug(fmt.Sprintf(format, a...))
	}
}

// Write an INFO message with key value pairs.
func (ml *MultiLogger) Info(msg any, kvs ...any) {
	for _, l := range ml.loggers {
		l.Info(msg, kvs...)
	}
}

// Write an INFO message with fmt.Sprintf.
func (ml *MultiLogger) Infof(format string, a ...any) {
	for _, l := range ml.loggers {
		l.Info(fmt.Sprintf(format, a...))
	}
}

// Write a WARN message with key value pairs.
func (ml *MultiLogger) Warn(msg any, kvs ...any) {
	for _, l := range ml.loggers {
		l.Warn(msg, kvs...)
	}
}

// Write a WARN message with fmt.Sprintf.
func (ml *MultiLogger) Warnf(format string, a ...any) {
	for _, l := range ml.loggers {
		l.Warn(fmt.Sprintf(format, a...))
	}
}

// Write an ERROR message with key value pairs.
func (ml *MultiLogger) Error(msg any, kvs ...any) {
	for _, l := range ml.loggers {
		l.Error(msg, kvs...)
	}
}

// Write an ERROR message with fmt.Sprintf.
func (ml *MultiLogger) Errorf(format string, a ...any) {
	for _, l := range ml.loggers {
		l.Error(fmt.Sprintf(format, a...))
	}
}

// Write a FATAL message with key value pairs.
func (ml *MultiLogger) Fatal(msg any, kvs ...any) {
	for _, l := range ml.loggers {
		l.Error(msg, kvs...)
	}

	os.Exit(1)
}

// Write a FATAL message with fmt.Sprintf.
func (ml *MultiLogger) Fatalf(format string, a ...any) {
	for _, l := range ml.loggers {
		l.Error(fmt.Sprintf(format, a...))
	}

	os.Exit(1)
}

// Write a PRINT message with key value pairs.
func (ml *MultiLogger) Print(msg any, kvs ...any) {
	for _, l := range ml.loggers {
		l.Print(msg, kvs...)
	}
}

// Write a PRINT message with fmt.Sprintf.
func (ml *MultiLogger) Printf(format string, a ...any) {
	for _, l := range ml.loggers {
		l.Print(fmt.Sprintf(format, a...))
	}
}
