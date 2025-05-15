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
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

// AttachConsole attaches the current process to an existing console.
func AttachConsole() error {
	const AttachParentProcess = ^uintptr(0)

	proc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("AttachConsole")

	r1, _, err := proc.Call(AttachParentProcess)
	if r1 == 0 {
		errno, ok := err.(syscall.Errno) //nolint:errorlint // reason: only used once.

		if ok && errno == 6 {
			return nil
		}

		return err
	}

	return nil
}

// AttachConsole attaches the current process to an existing console.
func AttachConsoleW() error { //nolint:gocyclo,cyclop // reason: not complex.
	wstdin, wstdout, wstderr := os.Stdin, os.Stdout, os.Stderr

	if wstdin == nil || wstdout == nil || wstderr == nil {
		return fmt.Errorf("attachconsole: %v %v %v", wstdin, wstdout, wstderr) //nolint:lll,err113 // reason: variables are required for error message.
	}

	stdin, _ := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	stdout, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	stderr, _ := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)

	var invalid syscall.Handle

	con := invalid

	if stdin == invalid || stdout == invalid || stderr == invalid {
		err := AttachConsole()
		if err != nil {
			return fmt.Errorf("attachconsole: %w", err)
		}

		if stdin == invalid {
			stdin, _ = syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
		}

		if stdout == invalid {
			stdout, _ = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
			con = stdout
		}

		if stderr == invalid {
			stderr, _ = syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)
			con = stderr
		}
	}

	if con != invalid {
		hand := windows.Handle(con)

		var mode uint32

		err := windows.GetConsoleMode(hand, &mode)
		if err != nil {
			return fmt.Errorf("GetConsoleMode: %w", err)
		}

		err = windows.SetConsoleMode(hand, mode&^windows.DISABLE_NEWLINE_AUTO_RETURN)
		if err != nil {
			return fmt.Errorf("SetConsoleMode: %w", err)
		}
	}

	if stdin != invalid {
		os.Stdin = os.NewFile(uintptr(stdin), "stdin")
	}

	if stdout != invalid {
		os.Stdout = os.NewFile(uintptr(stdout), "stdout")
	}

	if stderr != invalid {
		os.Stderr = os.NewFile(uintptr(stderr), "stderr")
	}

	return nil
}
