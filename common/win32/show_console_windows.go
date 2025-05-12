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

package win32

import (
	"syscall"
)

// ShowConsole shows or hides the console window.
func ShowConsole(show bool) error {
	getConsoleWindow := syscall.MustLoadDLL("kernel32.dll").MustFindProc("GetConsoleWindow")
	showWindow := syscall.MustLoadDLL("user32.dll").MustFindProc("ShowWindow")

	hwnd, _, err := getConsoleWindow.Call()
	if hwnd == 0 {
		return nil
	}

	if err != nil {
		return err
	}

	if show {
		var restore uintptr = 9

		_, _, err := showWindow.Call(hwnd, restore)
		if err != nil {
			return err
		}
	} else {
		var SwHide uintptr // 0

		_, _, err := showWindow.Call(hwnd, SwHide)
		if err != nil {
			return err
		}
	}

	return nil
}
