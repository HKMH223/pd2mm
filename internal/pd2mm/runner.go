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

package pd2mm

import (
	"github.com/hkmh223/pd2mm/common/benchmark"
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/io"
	"github.com/hkmh223/pd2mm/internal/lang"
)

var IsRunning bool //nolint:gochecknoglobals // allowed

// Run runs the program.
func (f Flags) Run(configs []Config, update func()) {
	errCh := make(chan error, 3) //nolint:mnd // allowed

	IsRunning = true

	go func() {
		logger.SharedLogger.Info(lang.Lang("startingNotify"))

		go f.RunWithError(configs, errCh)

		for err := range errCh {
			if err != nil {
				logger.SharedLogger.Errorf("%s %v", lang.Lang("errorNotify"), err)

				IsRunning = false
			}
		}

		IsRunning = false

		logger.SharedLogger.Info(lang.Lang("doneNotify"))

		update()
	}()
}

// RunWithError runs the program with error handling.
func (f Flags) RunWithError(configs []Config, errCh chan<- error) {
	defer close(errCh)

	for _, config := range configs {
		err := benchmark.Timer(func() error {
			f.runner(config)
			return nil
		}, "Start", func(methodName, elapsedTime string) {
			logger.SharedLogger.Infof("%s took %s", methodName, elapsedTime)
		})
		if err != nil {
			logger.SharedLogger.Fatal("Failed to benchmark", "err", err)
		}
	}
}

// runner starts the extraction and processing of mods.
func (f Flags) runner(config Config) {
	for _, search := range config.Mods {
		logger.SharedLogger.Info(lang.Lang("deleteNotify"), "path", search.Output)

		if err := filesystem.DeleteDirectory(filesystem.FromCwd(search.Output)); err != nil {
			logger.SharedLogger.Warn("Failed to delete directory", "path", search.Output, "err", err)
		}
	}

	for _, search := range config.Mods {
		if err := io.Extract(*f.Flags, search); err != nil {
			logger.SharedLogger.Error("Failed to extract mods", "err", err)
			continue
		}

		if err := config.Process(PathSearch{PathSearch: &search}); err != nil {
			logger.SharedLogger.Error("Failed to process mods", "err", err)
			continue
		}
	}
}
