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
	"os"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/otiai10/copy"
)

// CopyFile copies a file from the source to the destination.
// It skips files that are not allowed to be copied by pathCheck.
func CopyFile(src, dest string) error {
	return filesystem.Copy(src, dest, copy.Options{Skip: func(_ os.FileInfo, src, dest string) (bool, error) { //nolint:exhaustruct // allowed
		return PathCheck(src, dest), nil
	}, PermissionControl: copy.AddPermission(0o666)}) //nolint:mnd // allowed
}

// PathCheck checks if a file is allowed to be copied by the given source and destination paths.
func PathCheck(src, dest string) bool {
	srcResult, srcCheck := filesystem.CheckPathForProblemLocations(src)
	destResult, destCheck := filesystem.CheckPathForProblemLocations(dest)

	if srcResult || destResult {
		return !IsSafeAction(src, srcCheck) || !IsSafeAction(dest, destCheck)
	}

	return false
}

// Log based off of the PathCheck action type, returns a value based on whether to still copy files based off of the action.
func IsSafeAction(result string, check filesystem.PathCheck) bool {
	switch check.Action {
	case filesystem.PathCheckActionWarn:
		logger.SharedLogger.Warn("Problematic path located", "path", result, "type", check.Type, "target", check.Target)
		return true
	default:
		logger.SharedLogger.Error("Problematic path located", "path", result, "type", check.Type, "target", check.Target)
		return false
	}
}
