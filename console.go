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

package main

import (
	"io"
	"os"
	"slices"
	"sync"

	"github.com/hkmh223/pd2mm/common/benchmark"
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/lang"
	"github.com/hkmh223/pd2mm/internal/pd2mm"
)

var fileTypes = []string{".jsonc", ".json"} //nolint:gochecknoglobals // allowed

// StartApp is the main entry point for pd2mm.
func StartApp(logFile io.Writer) {
	logger.SharedLogger = logger.NewMultiLogger(logFile, os.Stdout)

	util.DrawWatermark([]string{lang.Lang("programName"), lang.Lang("watermarkPart1"), lang.Lang("watermarkPart2")}, func(s string) {
		logger.SharedLogger.Info(s)
	})

	if flags.Version {
		version()
		return
	}

	if util.IsFlagPassed("config") && flags.Config == "" {
		logger.SharedLogger.Fatal("Flag 'config' cannot be nil or empty")
	}

	var entries []string

	if util.IsFlagPassed("config") {
		entries = []string{flags.Config}
	} else {
		files, err := filesystem.GetTopFiles(lang.Lang("programName"))
		if err != nil {
			logger.SharedLogger.Fatal("Failed to get files", "err", err)
		}

		for _, file := range files {
			if slices.Contains(fileTypes, filesystem.GetFileExtension(file)) {
				entries = append(entries, filesystem.FromCwd(lang.Lang("programName"), file))
			}
		}
	}

	Run(entries)
}

// Run runs the program.
func Run(entries []string) {
	errCh := make(chan error, 3) //nolint:mnd // allowed

	var wait sync.WaitGroup

	wait.Add(1)

	go func() {
		defer wait.Done()

		logger.SharedLogger.Info(lang.Lang("startingNotify"))

		RunWithError(entries, errCh)

		for err := range errCh {
			if err != nil {
				logger.SharedLogger.Errorf("%s %v", lang.Lang("errorNotify"), err)
			}
		}

		logger.SharedLogger.Info(lang.Lang("doneNotify"))
	}()

	wait.Wait()
}

// RunWithError runs the program with error handling.
func RunWithError(entries []string, errCh chan<- error) {
	defer close(errCh)

	for _, entry := range entries {
		if c, err := pd2mm.Read(entry); err == nil {
			err := benchmark.Timer(func() error {
				flags.Start(c)
				return nil
			}, "Start", func(methodName, elapsedTime string) {
				logger.SharedLogger.Warnf("%s took %s", methodName, elapsedTime)
			})
			if err != nil {
				logger.SharedLogger.Fatal("Failed to benchmark", "err", err)
			}
		} else {
			errCh <- err

			return
		}
	}
}
