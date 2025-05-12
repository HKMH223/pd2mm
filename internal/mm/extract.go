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
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/sevenzip"
	"github.com/hkmh223/pd2mm/internal/lang"
)

func Extract(search PathSearch) error {
	if err := filesystem.DeleteDirectory(search.Extract); err != nil {
		logger.SharedLogger.Warn("Failed to delete directory", "path", search.Extract, "err", err)
	}

	logger.SharedLogger.Info(lang.Lang("extracting"), "source", search.Path, "destination", search.Extract)

	return extract(search.Path, search.Extract)
}

func extract(src string, dst string) error {
	files := filesystem.GetFiles(src)

	for _, file := range files {
		if _, err := sevenzip.Extract(file, dst, false); err != nil {
			return err
		}
	}

	return nil
}
