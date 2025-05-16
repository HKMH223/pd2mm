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
	"errors"
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

var ErrFreeConsole = errors.New("FreeConsole returned 0")

// FreeConsole frees the console associated with the current process.
func FreeConsole() error {
	proc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("FreeConsole")

	r, _, e := proc.Call()
	if int32(r) == 0 {
		if e != nil && !errors.Is(e, windows.ERROR_SUCCESS) {
			return fmt.Errorf("FreeConsole failed: %w", e)
		}

		return ErrFreeConsole
	}

	return nil
}
