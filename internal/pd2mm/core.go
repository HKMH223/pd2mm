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
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hkmh223/pd2mm/common/errors"
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/safe"
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/io"
	"github.com/hkmh223/pd2mm/internal/lang"
)

type MError = errors.MError

// Process handles copying files with the given PathSearch.
func (c Config) Process(ps PathSearch) error {
	if err := c.process(ps); err != nil {
		return err
	}

	if err := c.copyAdditional(ps); err != nil {
		return err
	}

	return nil
}

// process handles copying files with the given PathSearch.
func (c Config) process(search PathSearch) error {
	directories, err := filesystem.GetTopDirectories(filesystem.FromCwd(search.Extract.Path))
	if err != nil {
		logger.SharedLogger.Warn("failed to get directories", "path", search.Extract.Path, "err", err)
		return err
	}

	for _, directory := range directories {
		c.checkIncludeData(filesystem.Normalize(filepath.Join(search.Extract.Path, directory)), search)
	}

	if search.Export.Path != "" {
		if err := io.CopyFile(search.Output.Path, search.Export.Path); err != nil {
			return &MError{Header: "process", Message: fmt.Sprintf("failed to copy '%s' to '%s'", search.Output.Path, search.Export.Path), Err: err}
		}
	}

	return nil
}

// Handle checkIncludeData settings for a given path.
func (c Config) checkIncludeData(path string, search PathSearch) {
	files := filesystem.GetFiles(filesystem.FromCwd(path))

	for _, file := range files {
		source := filesystem.Normalize(file)

		for _, include := range search.Include {
			if !strings.Contains(source, search.FormatString(include.Path)) {
				continue
			}

			if err := c.copyExpected(source, search.FormatString(include.To), false, search); err != nil {
				logger.SharedLogger.Error("failed to copy", "source", source, "destination", search.FormatString(include.To), "err", err)
			}
		}

		if c.checkExcludeData(source, search) {
			break
		}
	}
}

// Handle checkExcludeData settings for a given path.
func (c Config) checkExcludeData(path string, search PathSearch) bool {
	parts := strings.Split(filesystem.Normalize(path), "/")

	for _, exclude := range search.Exclude {
		if util.ContainsSubslice(parts, search.FormatSlice(exclude)) {
			return false
		}
	}

	skip, err := c.checkExpectsData(path, search)
	if err != nil {
		logger.SharedLogger.Error("failed to copy expected paths", "err", err)
		return true
	}

	return skip
}

// Handle expected settings for a given path.
func (c Config) checkExpectsData(path string, search PathSearch) (bool, error) {
	source := strings.Split(path, "/")

	for _, expect := range search.Expects {
		if len(source) < len(expect.Path) {
			continue
		}

		if slices.Contains(expect.Path, filesystem.GetFileExtension(safe.Slice(source, len(source)-1))) {
			return c.expectedIsFile(source, search, expect)
		} else if util.Matches(source, expect.Path) == len(expect.Path) {
			return c.expectedIsDirectory(source, search, expect)
		}
	}

	return false, nil
}

// Handle expected data as a file.
func (c Config) expectedIsFile(source []string, search PathSearch, expect data.Expect) (bool, error) {
	path := strings.Join(source, "/")

	destination := filepath.Join(search.Output.Path, fixDestination(source, search, expect, false))
	if destination == "" {
		return false, nil
	}

	destination = filesystem.Normalize(filesystem.FromCwd(destination))

	if util.ContainsSubslice(source, expect.Require) && util.ContainsSubslice(strings.Split(destination, "/"), expect.Require) {
		path = strings.Join(safe.Range(source, 0, len(source)-len(expect.Require)), "/")

		dest := strings.Split(destination, "/")
		destination = strings.Join(safe.Range(dest, 0, len(dest)-len(expect.Require)), "/")
	}

	if err := c.copyExpected(path, destination, false, search); err != nil {
		return false, &MError{Header: "expects", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", path, destination), Err: err}
	}

	return true, nil
}

// Handle expected data as a directory.
func (c Config) expectedIsDirectory(source []string, search PathSearch, expect data.Expect) (bool, error) {
	destination := fixDestination(source, search, expect, true)
	if destination == "" {
		return false, nil
	}

	destination = filesystem.FromCwd(destination)

	index := safe.HasIndex(source, safe.Slice(expect.Path, 0))
	src := strings.Join(safe.Range(source, 0, index), "/")

	if expect.Exclusive {
		src = strings.Join(safe.Range(source, 0, index+1), "/")
	}

	if err := c.copyExpected(src, destination, false, search); err != nil {
		return false, &MError{Header: "expectedIsDirectory", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, destination), Err: err}
	}

	return true, nil
}

// Handle non-contextual file copyAdditional.
//
//nolint:lll // reason: logging.
func (c Config) copyAdditional(search PathSearch) error {
	for _, copy := range search.Copy {
		logger.SharedLogger.Info(lang.Lang("copyingNotify"), "source", search.FormatString(copy.From), "destination", search.FormatString(copy.To))

		src, dest := search.FormatString(copy.From), search.FormatString(copy.To)
		if err := io.CopyFile(src, dest); err != nil {
			return &MError{Header: "copyAdditional", Message: fmt.Sprintf("failed to copy '%s' to '%s'", src, dest), Err: err}
		}
	}

	return nil
}

// fixDestination fixes the destination path based on the provided PathSearch and Expect data.
func fixDestination(parts []string, search PathSearch, expect data.Expect, dir bool) string {
	result := strings.Join(parts, "/")

	if dir {
		index := safe.HasIndex(parts, safe.Slice(expect.Path, 0))
		result = strings.Join(safe.Range(parts, 0, index), "/")

		if expect.Exclusive {
			result = strings.Join(safe.Range(parts, 0, index+1), "/")
		}
	}

	results := strings.Split(result, "/")
	expect.Require = append(expect.Require, safe.Slice(results, len(results)-1))
	base := strings.Join(expect.Require, "/")

	if dir {
		base = strings.Join(safe.Range(expect.Require, 0, len(expect.Require)-expect.Base), "/")
		return filepath.Join(search.Output.Path, base)
	}

	return base
}

// copyExpected copies the source file or directory to the destination based on the provided PathSearch.
func (c Config) copyExpected(src, dest string, expected bool, search PathSearch) error {
	src = filesystem.Normalize(src)
	dest = filesystem.Normalize(dest)

	for _, rename := range search.Rename {
		pathFmt := search.FormatSlice(rename.Path)
		fromFmt, toFmt := search.FormatSlice(rename.From), search.FormatSlice(rename.To)
		parts := strings.Split(src, "/")

		if util.ContainsSubslice(parts, pathFmt) {
			result := util.ReplaceSubslice(strings.Split(dest, "/"), fromFmt, toFmt)
			dest = strings.Join(result, "/")
		}
	}

	if expected {
		if err := c.parseExpectedAndCopy(src, dest); err != nil {
			return &MError{Header: "copyExpected", Message: fmt.Sprintf("failed to copy '%s' to '%s'", src, dest), Err: err}
		}
	}

	logger.SharedLogger.Info(lang.Lang("copyingNotify"), "source", src, "destination", dest)

	if err := io.CopyFile(src, dest); err != nil {
		return &MError{Header: "copyExpected", Message: fmt.Sprintf("failed to copy '%s' to '%s'", src, dest), Err: err}
	}

	return nil
}

// Parse expected file paths and copy them.
func (c Config) parseExpectedAndCopy(src, dest string) error {
	parts := strings.Split(src, "/")
	source := safe.Range(parts, 0, len(parts)-1)

	// Combine the normalized destination with the source directory name
	src, dest = strings.Join(source, "/"), filepath.Join(dest, safe.Slice(source, len(source)-1))
	logger.SharedLogger.Info(lang.Lang("copyingNotify"), "source", src, "destination", dest)

	if err := io.CopyFile(src, dest); err != nil {
		return &MError{Header: "parseExpectedAndCopy", Message: fmt.Sprintf("failed to copy '%s' to '%s'", src, dest), Err: err}
	}

	return nil
}
