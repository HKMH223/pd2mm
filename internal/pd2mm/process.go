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
	directories, err := filesystem.GetTopDirectories(filesystem.FromCwd(search.Extract))
	if err != nil {
		logger.SharedLogger.Debug("Failed to get directories", "path", search.Extract, "err", err)
		return err
	}

	for _, directory := range directories {
		c.checkIncludeData(strings.ReplaceAll(filepath.Join(search.Extract, directory), "\\", "/"), search)
	}

	if search.Export != "" {
		if err := copyFile(search.Output, search.Export); err != nil {
			return &MError{Header: "process", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", search.Output, search.Export), Err: err}
		}
	}

	return nil
}

// Handle checkIncludeData settings for a given path.
func (c Config) checkIncludeData(path string, search PathSearch) {
	files := filesystem.GetFiles(filesystem.FromCwd(path))

	for _, file := range files {
		source := strings.ReplaceAll(file, "\\", "/")

		for _, include := range search.Include {
			if !strings.Contains(source, search.formatString(include.Path)) {
				continue
			}

			if err := c.copyExpected(source, search.formatString(include.To), false, search); err != nil {
				logger.SharedLogger.Error("Failed to copy", "source", source, "destination", search.formatString(include.To), "err", err)
			}
		}

		if c.checkExcludeData(source, search) {
			break
		}
	}
}

// Handle checkExcludeData settings for a given path.
func (c Config) checkExcludeData(path string, search PathSearch) bool {
	parts := strings.Split(strings.ReplaceAll(path, "\\", "/"), "/")

	for _, exclude := range search.Exclude {
		if util.ContainsSubslice(parts, search.formatSlice(exclude)) {
			return false
		}
	}

	skip, err := c.checkExpectsData(path, search)
	if err != nil {
		logger.SharedLogger.Fatal("Failed to copy expected paths", "err", err)
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
func (c Config) expectedIsFile(source []string, search PathSearch, expect Expect) (bool, error) {
	path := strings.Join(source, "/")

	destination := filepath.Join(search.Output, fixDestination(source, search, expect, false))
	if destination == "" {
		return false, nil
	}

	destination = strings.ReplaceAll(filesystem.FromCwd(destination), "\\", "/")

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
func (c Config) expectedIsDirectory(source []string, search PathSearch, expect Expect) (bool, error) {
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
		return false, &MError{Header: "expects", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, destination), Err: err}
	}

	return true, nil
}

// Handle non-contextual file copyAdditional.
//
//nolint:lll // allowed
func (c Config) copyAdditional(search PathSearch) error {
	for _, copy := range search.Copy {
		logger.SharedLogger.Info(lang.Lang("copyingNotify"), "source", search.formatString(copy.From), "destination", search.formatString(copy.To))

		src, dest := search.formatString(copy.From), search.formatString(copy.To)
		if err := copyFile(src, dest); err != nil {
			return &MError{Header: "copies", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, dest), Err: err}
		}
	}

	return nil
}

// fixDestination fixes the destination path based on the provided PathSearch and Expect data.
func fixDestination(parts []string, search PathSearch, expect Expect, dir bool) string {
	result := strings.Join(parts, "/")

	if dir {
		index := safe.HasIndex(parts, safe.Slice(expect.Path, 0))
		result = strings.Join(safe.Range(parts, 0, index), "/")

		if expect.Exclusive {
			result = strings.Join(safe.Range(parts, 0, index+1), "/")
		}
	}

	results := strings.Split(result, "/")
	destination := append(expect.Require, safe.Slice(results, len(results)-1)) //nolint:gocritic // allowed
	base := strings.Join(destination, "/")

	if dir {
		base = strings.Join(safe.Range(destination, 0, len(destination)-expect.Base), "/")
		return filepath.Join(search.Output, base)
	}

	return base
}

// copyExpected copies the source file or directory to the destination based on the provided PathSearch.
func (c Config) copyExpected(src, dest string, expected bool, search PathSearch) error {
	src = strings.ReplaceAll(src, "\\", "/")
	dest = strings.ReplaceAll(dest, "\\", "/")

	for _, rename := range search.Rename {
		pathFmt := search.formatSlice(rename.Path)
		fromFmt, toFmt := search.formatSlice(rename.From), search.formatSlice(rename.To)
		parts := strings.Split(src, "/")

		if util.ContainsSubslice(parts, pathFmt) {
			result := util.ReplaceSubslice(strings.Split(dest, "/"), fromFmt, toFmt)
			dest = strings.Join(result, "/")
		}
	}

	if expected {
		if err := c.parseExpectedAndCopy(src, dest); err != nil {
			return &MError{Header: "copyExpected", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, dest), Err: err}
		}
	}

	logger.SharedLogger.Info(lang.Lang("copyingNotify"), "source", src, "destination", dest)

	if err := copyFile(src, dest); err != nil {
		return &MError{Header: "copyExpected", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, dest), Err: err}
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

	if err := copyFile(src, dest); err != nil {
		return &MError{Header: "parseExpectedAndCopy", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, dest), Err: err}
	}

	return nil
}

// Format a slice of paths using the current PathSearch settings.
func (ps PathSearch) formatSlice(slice []string) []string {
	result := []string{}

	for _, item := range slice {
		result = append(result, ps.formatString(item))
	}

	return result
}

// Replace keywords with relevant PathSearch settings.
func (ps PathSearch) formatString(str string) string {
	return util.Format(str, map[string]string{
		"{path}":    ps.Path,
		"{output}":  ps.Output,
		"{extract}": ps.Extract,
		"{export}":  ps.Export,
	})
}
