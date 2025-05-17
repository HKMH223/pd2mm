//go:build windows

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

package filesystem

import (
	"syscall"
)

// clearReadOnlyAttr clears the read-only attribute of a file or directory.
func clearReadOnlyAttr(path string) error {
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	attrs, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return err
	}

	if attrs&syscall.FILE_ATTRIBUTE_READONLY != 0 {
		attrs &^= syscall.FILE_ATTRIBUTE_READONLY

		err = syscall.SetFileAttributes(pointer, attrs)
		if err != nil {
			return err
		}
	}

	return nil
}
