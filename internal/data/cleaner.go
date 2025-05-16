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

	"github.com/hkmh223/pd2mm/common/errors"
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/lang"
)

// Delete all files in the Extract path that are not in the list of exclusions.
func (ps PathSearch) CleanExtractDirectory() error {
	target := filesystem.FromCwd(ps.Extract.Path)

	logger.SharedLogger.Info(lang.Lang("deleteNotify"), "path", target)

	if err := filesystem.DeleteDirectory(target, func(s string) bool {
		return skip(s, ps, ps.Extract)
	}); err != nil {
		return &errors.MError{Header: "CleanExtract", Message: "failed to delete directory: " + target, Err: err}
	}

	if err := filesystem.DeleteEmptyDirectories(target); err != nil {
		return err
	}

	return nil
}

// Delete all files in the Export path that are not in the list of exclusions.
func (ps PathSearch) CleanExportDirectory() error {
	target := filesystem.FromCwd(ps.Export.Path)

	logger.SharedLogger.Info(lang.Lang("deleteNotify"), "path", target)

	if err := filesystem.DeleteDirectory(target, func(s string) bool {
		return skip(s, ps, ps.Export)
	}); err != nil {
		return &errors.MError{Header: "CleanExport", Message: "failed to delete directory: " + target, Err: err}
	}

	if err := filesystem.DeleteEmptyDirectories(target); err != nil {
		return err
	}

	return nil
}

// Delete all files in the Output path that are not in the list of exclusions.
func (ps PathSearch) CleanOutputDirectory() error {
	target := filesystem.FromCwd(ps.Output.Path)

	logger.SharedLogger.Info(lang.Lang("deleteNotify"), "path", target)

	if err := filesystem.DeleteDirectory(target, func(s string) bool {
		return skip(s, ps, ps.Output)
	}); err != nil {
		return &errors.MError{Header: "CleanOutput", Message: "failed to delete directory: " + target, Err: err}
	}

	if err := filesystem.DeleteEmptyDirectories(target); err != nil {
		return err
	}

	return nil
}

// Check if the file name should be excluded.
func skip(name string, ps PathSearch, pi PathInfo) bool {
	normalized := strings.Split(filesystem.Normalize(name), "/")

	for _, exclude := range pi.ExcludeClean {
		excludeNormalized := ps.formatSubSlices(exclude)

		logger.SharedLogger.Info(normalized)
		logger.SharedLogger.Info(excludeNormalized)

		if util.ContainsSubslice(normalized, excludeNormalized) {
			return true
		}
	}

	return false
}

// Format a slice of file names to be used in the exclude function.
func (ps PathSearch) formatSubSlices(slice []string) []string {
	normalized := ps.FormatSlice(filesystem.NormalizeSlice(slice))
	result := []string{}

	for _, str := range normalized {
		result = slices.Concat(result, strings.Split(str, "/"))
	}

	return result
}
