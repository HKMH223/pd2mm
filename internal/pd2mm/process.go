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
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/lang"
)

type MError = errors.MError

// Process handles copying files with the given PathSearch.
func (c Config) Process(ps PathSearch) error {
	if err := c.process(ps); err != nil {
		return err
	}

	if err := c.copies(ps); err != nil {
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
		c.include(strings.ReplaceAll(filepath.Join(search.Extract, directory), "\\", "/"), search)
	}

	if search.Export != "" {
		if err := copyFile(search.Output, search.Export); err != nil {
			return &MError{Header: "process", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", search.Output, search.Export), Err: err}
		}
	}

	return nil
}

// Handle include settings for a given path.
func (c Config) include(path string, search PathSearch) {
	files := filesystem.GetFiles(filesystem.FromCwd(path))

	for _, file := range files {
		source := strings.ReplaceAll(file, "\\", "/")

		for _, include := range search.Include {
			if !strings.Contains(source, search.format(include.Path)) {
				continue
			}

			if err := c.copyExpected(source, search.format(include.To), false, search); err != nil {
				logger.SharedLogger.Error("Failed to copy", "source", source, "destination", search.format(include.To), "err", err)
			}
		}

		if c.exclude(source, search) {
			break
		}
	}
}

// Handle exclude settings for a given path.
func (c Config) exclude(path string, search PathSearch) bool {
	for _, exclude := range search.Exclude {
		if strings.Contains(path, search.format(exclude)) {
			return false
		}
	}

	skip, err := c.expects(path, search)
	if err != nil {
		logger.SharedLogger.Fatal("Failed to copy expected paths", "err", err)
	}

	return skip
}

// Handle expected settings for a given path.
func (c Config) expects(path string, search PathSearch) (bool, error) {
	var destination string

	source := strings.Split(path, "/")

	for _, expect := range search.Expects {
		if len(source) < len(expect.Path) {
			continue
		}

		if slices.Contains(expect.Path, filesystem.GetFileExtension(source[len(source)-1])) { //nolint:nestif // allowed
			destination = filepath.Join(search.Output, fixDestination(source, search, expect, false))
			if destination == "" {
				return false, nil
			}

			if err := c.copyExpected(path, destination, false, search); err != nil {
				return false, &MError{Header: "expects", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", path, destination), Err: err}
			}

			return true, nil
		} else if util.Matches(source, expect.Path) == len(expect.Path) {
			destination = fixDestination(source, search, expect, true)
			if destination == "" {
				return false, nil
			}

			index := slices.Index(source, expect.Path[0])
			if index == -1 {
				logger.SharedLogger.Fatalf("%v expected %s but it was not found (index -1)", source, expect.Path[0])
			}

			src := strings.Join(source[:index], "/")
			if err := c.copyExpected(src, destination, false, search); err != nil {
				return false, &MError{Header: "expects", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, destination), Err: err}
			} // path

			return true, nil
		}
	}

	return false, nil
}

// Handle non-contextual file copies.
func (c Config) copies(search PathSearch) error {
	for _, copy := range search.Copy {
		logger.SharedLogger.Info(lang.Lang("copying"), "source", search.format(copy.From), "destination", search.format(copy.To))

		src, dest := search.format(copy.From), search.format(copy.To)
		if err := copyFile(src, dest); err != nil {
			return &MError{Header: "copies", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, dest), Err: err}
		}
	}

	return nil
}

// fixDestination fixes the destination path based on the provided PathSearch and Expect data.
func fixDestination(parts []string, search PathSearch, expect Expect, dir bool) string {
	// Join initial segments minus the length of expected path.
	result := strings.Join(parts, "/")

	if dir {
		index := slices.Index(parts, expect.Path[0])
		if index == -1 {
			logger.SharedLogger.Fatalf("%v expected %s but it was not found (index -1)", parts, expect.Path[0])
		}

		result = strings.Join(parts[:index], "/") // parts[:len(parts)-len(ex.Path)]
	}

	results := strings.Split(result, "/")

	// Append last element of newParts to the expected requirements.
	destination := append(expect.Require, results[len(results)-1]) //nolint:gocritic // allowed

	// Join finalParts minus the length of expected base path.
	base := strings.Join(destination, "/")
	if dir {
		base = strings.Join(destination[:len(destination)-expect.Base], "/")
		return filesystem.FromCwd(filepath.Join(search.Output, base))
	}

	return filesystem.FromCwd(base)
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

	logger.SharedLogger.Info(lang.Lang("copying"), "source", src, "destination", dest)

	if err := copyFile(src, dest); err != nil {
		return &MError{Header: "copyExpected", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, dest), Err: err}
	}

	return nil
}

// Parse expected file paths and copy them.
func (c Config) parseExpectedAndCopy(src, dest string) error {
	parts := strings.Split(src, "/")
	source := parts[:len(parts)-1]

	// Combine the normalized destination with the source directory name
	src, dest = strings.Join(source, "/"), filepath.Join(dest, source[len(source)-1])
	logger.SharedLogger.Info(lang.Lang("copying"), "source", src, "destination", dest)

	if err := copyFile(src, dest); err != nil {
		return &MError{Header: "parseExpectedAndCopy", Message: fmt.Sprintf("Failed to copy '%s' to '%s'", src, dest), Err: err}
	}

	return nil
}

// Format a slice of paths using the current PathSearch settings.
func (ps PathSearch) formatSlice(slice []string) []string {
	result := []string{}

	for _, item := range slice {
		result = append(result, ps.format(item))
	}

	return result
}

// Replace keywords with relevant PathSearch settings.
func (ps PathSearch) format(str string) string {
	return util.Format(str, map[string]string{
		"{path}":    ps.Path,
		"{output}":  ps.Output,
		"{extract}": ps.Extract,
		"{export}":  ps.Export,
	})
}
