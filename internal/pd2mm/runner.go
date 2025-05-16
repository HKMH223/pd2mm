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
	"sync"
	"sync/atomic"

	"github.com/hkmh223/pd2mm/common/benchmark"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/io"
	"github.com/hkmh223/pd2mm/internal/lang"
)

var SharedRunner = NewRunner(func() error { return nil }) //nolint:gochecknoglobals // reason: used by generic cleaning functions

type Runner struct {
	mu sync.Mutex

	isRunning atomic.Bool
	Update    func() error
}

// NewRunner creates a new Runner.
func NewRunner(update func() error) *Runner {
	return &Runner{ //nolint:exhaustruct // reason: value is set
		Update: update,
	}
}

func (r *Runner) RegisterUpdate(update func() error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Update = update
}

// IsActive returns true if the Runner is currently cleaning.
func (r *Runner) IsActive() bool {
	return r.isRunning.Load()
}

// Run runs the program.
func (r *Runner) Run(flags Flags, configs []Config) {
	r.isRunning.Store(true)

	defer func() {
		r.isRunning.Store(false)
		logger.SharedLogger.Info(lang.Lang("doneRunnerNotify"))

		if err := r.Update(); err != nil {
			logger.SharedLogger.Error(err)
		}
	}()

	logger.SharedLogger.Info(lang.Lang("startingRunnerNotify"))

	errCh := make(chan error, 1)
	go flags.RunWithError(configs, errCh)

	for err := range errCh {
		if err != nil {
			logger.SharedLogger.Errorf("%s %v", lang.Lang("errorNotify"), err)
		}
	}
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
			errCh <- err

			return
		}
	}
}

// runner starts the extraction and processing of mods.
func (f Flags) runner(config Config) {
	SharedCleaner.Clean([]Config{config}, Output, func() error {
		for _, search := range config.Mods {
			io.Extract(*f.Flags, search)

			if err := config.Process(PathSearch{PathSearch: &search}); err != nil {
				logger.SharedLogger.Error("failed to process mods", "err", err)
				continue
			}
		}

		return nil
	})
}
