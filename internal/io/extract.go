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

package io

import (
	"runtime"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/sevenzip"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/lang"
)

// Extract extracts the contents of an archive to a specified directory.
func Extract(flags data.Flags, search data.PathSearch) error {
	if err := filesystem.DeleteDirectory(filesystem.FromCwd(search.Extract)); err != nil {
		logger.SharedLogger.Warn("Failed to delete directory", "path", search.Extract, "err", err)
	}

	source := filesystem.FromCwd(search.Path)
	destination := filesystem.FromCwd(search.Extract)
	logger.SharedLogger.Info(lang.Lang("extractingNotify"), "source", source, "destination", destination)

	return extract(flags, source, destination)
}

// extract extracts the contents of an archive to a specified directory.
func extract(flags data.Flags, src, dest string) error {
	files := filesystem.GetFiles(src)

	bin := filesystem.Combine(flags.Bin, sevenzip.LinuxName)
	if runtime.GOOS == "windows" {
		bin = filesystem.Combine(flags.Bin, sevenzip.WindowsName)
	}

	for _, file := range files {
		if filesystem.Exists(bin) {
			if _, err := sevenzip.ExtractWithBin(file, dest, bin, false); err != nil {
				return err
			}
		} else {
			if _, err := sevenzip.Extract(file, dest, false, sevenzip.ExtractionOptions{HideWindow: true, Relative: false}); err != nil {
				return err
			}
		}
	}

	return nil
}
