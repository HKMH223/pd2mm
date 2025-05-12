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
	"os"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/otiai10/copy"
)

// copyFile copies a file from the source to the destination.
// It skips files that are not allowed to be copied by pathCheck.
func copyFile(src, dest string) error {
	return filesystem.Copy(src, dest, copy.Options{Skip: func(_ os.FileInfo, src, dest string) (bool, error) { //nolint:exhaustruct // allowed
		return pathCheck(src, dest), nil
	}})
}

// pathCheck checks if a file is allowed to be copied by the given source and destination paths.
func pathCheck(src, dest string) bool {
	if srcResult, srcCheck := filesystem.CheckPathForProblemLocations(src); srcResult {
		if srcCheck.Action == filesystem.PathCheckActionWarn {
			logger.SharedLogger.Warn("Problematic path located", "path", src, "type", srcCheck.Type, "target", srcCheck.Target)
			return false
		}

		logger.SharedLogger.Error("Problematic path located", "path", src, "type", srcCheck.Type, "target", srcCheck.Target)

		return true
	}

	if destResult, destCheck := filesystem.CheckPathForProblemLocations(dest); destResult {
		if destCheck.Action == filesystem.PathCheckActionWarn {
			logger.SharedLogger.Error("Problematic path located", "path", dest, "type", destCheck.Type, "target", destCheck.Target)
			return false
		}

		logger.SharedLogger.Error("Problematic path located", "path", dest, "type", destCheck.Type, "target", destCheck.Target)

		return true
	}

	return false
}
