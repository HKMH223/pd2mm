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
	"github.com/hkmh223/pd2mm/common/process"
	"golang.org/x/sys/windows"
)

// HideConsoleWindow frees the console associated with the current process.
func HideConsoleWindow() {
	if process.ProcGetWindowThreadProcessID.Find() != nil {
		return
	}

	pid := windows.GetCurrentProcessId()

	var cpid uint32
	if _, err := windows.GetWindowThreadProcessId(windows.HWND(GetConsoleWindow()), &cpid); err != nil {
		return
	}

	if pid == cpid {
		_ = FreeConsole()
	}
}
