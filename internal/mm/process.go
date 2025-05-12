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

package mm

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/lang"
)

func (c Config) Process(ps PathSearch) error {
	if err := c.process(ps); err != nil {
		return err
	}

	if err := c.copies(ps); err != nil {
		return err
	}

	return nil
}

func (c Config) process(search PathSearch) error {
	directories, err := filesystem.GetTopDirectories(search.Extract)
	if err != nil {
		logger.SharedLogger.Debug("Failed to get directories", "path", search.Extract, "err", err)
		return err
	}

	for _, directory := range directories {
		c.include(strings.ReplaceAll(filepath.Join(search.Extract, directory), "\\", "/"), search)
	}

	return nil
}

func (c Config) include(path string, search PathSearch) {
	files := filesystem.GetFiles(path)

	for _, file := range files {
		source := strings.ReplaceAll(file, "\\", "/")

		for _, include := range search.Include {
			if strings.Contains(source, search.format(include.Path)) {
				if err := c.copy(source, search.format(include.To), false, search); err != nil {
					logger.SharedLogger.Error("Failed to copy", "source", source, "destination", search.format(include.To), "err", err)
				}
			}
		}

		if c.exclude(source, search) {
			break
		}
	}
}

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

func (c Config) expects(path string, search PathSearch) (bool, error) {
	var destination string

	source := strings.Split(path, "/")

	for _, expect := range search.Expects {
		if len(source) < len(expect.Path) {
			continue
		}

		if slices.Contains(expect.Path, filesystem.GetFileExtension(source[len(source)-1])) { //nolint:nestif // allowed
			destination = filepath.Join(search.Export, fixDestination(source, search, expect, false))

			if destination != "" {
				if err := c.copy(path, destination, false, search); err != nil {
					return false, err
				}

				return true, nil
			}
		} else if util.Matches(source, expect.Path) == len(expect.Path) {
			destination = fixDestination(source, search, expect, true)

			if destination != "" {
				index := slices.Index(source, expect.Path[0])
				if index == -1 {
					logger.SharedLogger.Fatalf("%v expected %s but it was not found (index -1)", source, expect.Path[0])
				}

				if err := c.copy(strings.Join(source[:index], "/"), destination, false, search); err != nil {
					return false, err
				} // path

				return true, nil
			}
		}
	}

	return false, nil
}

func (c Config) copies(search PathSearch) error {
	for _, copy := range search.Copy {
		logger.SharedLogger.Info(lang.Lang("copying"), "source", search.format(copy.From), "destination", search.format(copy.To))

		if err := filesystem.Copy(search.format(copy.From), search.format(copy.To)); err != nil {
			return err
		}
	}

	return nil
}

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
		return filepath.Join(search.Export, base)
	}

	return base
}

func (c Config) copy(src, dst string, expected bool, search PathSearch) error {
	src = strings.ReplaceAll(src, "\\", "/")
	dst = strings.ReplaceAll(dst, "\\", "/")

	for _, rename := range search.Rename {
		pathFmt := search.formatSlice(rename.Path)
		fromFmt := search.formatSlice(rename.From)
		toFmt := search.formatSlice(rename.To)
		parts := strings.Split(src, "/")

		if util.ContainsSubslice(parts, pathFmt) {
			result := util.ReplaceSubslice(strings.Split(dst, "/"), fromFmt, toFmt)
			dst = strings.Join(result, "/")
		}
	}

	if expected {
		return c.copyExpected(src, dst)
	}

	logger.SharedLogger.Info(lang.Lang("copying"), "source", src, "destination", dst)

	return filesystem.Copy(src, dst)
}

func (c Config) copyExpected(src, dst string) error {
	parts := strings.Split(src, "/")
	source := parts[:len(parts)-1]

	// Combine the normalized destination with the source directory name
	src = strings.Join(source, "/")
	dst = filepath.Join(dst, source[len(source)-1])
	logger.SharedLogger.Info(lang.Lang("copying"), "source", src, "destination", dst)

	return filesystem.Copy(src, dst)
}

func (ps PathSearch) formatSlice(slice []string) []string {
	result := []string{}

	for _, item := range slice {
		result = append(result, ps.format(item))
	}

	return result
}

func (ps PathSearch) format(str string) string {
	format := strings.ReplaceAll(str, "{path}", ps.Path)
	format = strings.ReplaceAll(format, "{extract}", ps.Extract)
	format = strings.ReplaceAll(format, "{export}", ps.Export)

	return format
}
