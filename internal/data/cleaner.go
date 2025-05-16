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

package data

import (
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hkmh223/pd2mm/common/errors"
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/lang"
)

var SharedCleaner = NewCleaner(func() error { return nil }) //nolint:gochecknoglobals // reason: used by generic cleaning functions

type Cleaner struct {
	mu sync.Mutex

	isCleaning atomic.Bool
	Update     func() error
}

// NewCleaner creates a new cleaner.
func NewCleaner(update func() error) *Cleaner {
	return &Cleaner{ //nolint:exhaustruct // reason: value is set
		Update: update,
	}
}

// RegisterUpdate registers an update function that is called after the cleaner finishes cleaning.
func (c *Cleaner) RegisterUpdate(update func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Update = update
}

// IsActive returns true if the cleaner is currently cleaning.
func (c *Cleaner) IsActive() bool {
	return c.isCleaning.Load()
}

// Delete all files in the path that are not in the list of exclusions.
func (c *Cleaner) Clean(search PathSearch, info PathInfo) {
	c.isCleaning.Store(true)

	defer func() {
		c.isCleaning.Store(false)
		logger.SharedLogger.Info(lang.Lang("doneCleanerNotify"))

		if err := c.Update(); err != nil {
			logger.SharedLogger.Error(err)
		}
	}()

	logger.SharedLogger.Info(lang.Lang("startingCleanerNotify"))

	errCh := make(chan error, 1)
	go search.CleanWithError(info, errCh)

	for err := range errCh {
		if err != nil {
			logger.SharedLogger.Errorf("%s %v", lang.Lang("errorNotify"), err)
		}
	}
}

// Delete all files in the path that are not in the list of exclusions.
func (search PathSearch) CleanWithError(info PathInfo, errCh chan<- error) {
	defer close(errCh)

	target, err := filesystem.FromCwd(info.Path)
	if err != nil {
		errCh <- err

		return
	}

	logger.SharedLogger.Info(lang.Lang("deleteNotify"), "path", target)

	if err := filesystem.DeleteDirectory(target, func(s string) bool {
		return skip(s, search, info)
	}); err != nil {
		errCh <- &errors.MError{Header: "CleanWithError", Message: "failed to delete directory: " + target, Err: err}

		return
	}

	if err := filesystem.DeleteEmptyDirectories(target); err != nil {
		errCh <- err

		return
	}
}

// Check if the file name should be excluded.
func skip(name string, search PathSearch, info PathInfo) bool {
	normalized := strings.Split(filesystem.Normalize(name), "/")

	for _, exclude := range info.ExcludeClean {
		excludeNormalized := search.formatSubSlices(exclude)

		logger.SharedLogger.Info(normalized)
		logger.SharedLogger.Info(excludeNormalized)

		if util.ContainsSubslice(normalized, excludeNormalized) {
			return true
		}
	}

	return false
}

// Format a slice of file names to be used in the exclude function.
func (search PathSearch) formatSubSlices(slice []string) []string {
	normalized := search.FormatSlice(filesystem.NormalizeSlice(slice))

	var result []string

	for _, str := range normalized {
		result = slices.Concat(result, strings.Split(str, "/"))
	}

	return result
}
