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
	"io"
	"os"
)

// WindowConsoleHandle optionally allocated a console to a window process.
func WindowConsoleHandle(args []string, minArgs int, console, window func(in, out, err io.Writer), isConsole bool) error {
	if len(args) > minArgs { //nolint:nestif // reason: not complex
		if err := AttachConsoleW(); err != nil {
			return err
		}

		console(os.Stdin, os.Stdout, os.Stderr)
	} else {
		var err error

		var cIn, cOut, cErr io.Writer

		if isConsole {
			cIn, cOut, cErr, err = AllocConsole()
			if err != nil {
				return err
			}
		} else {
			cIn, cOut, cErr = os.Stdin, os.Stdout, os.Stderr
		}

		window(cIn, cOut, cErr)
	}

	return nil
}
