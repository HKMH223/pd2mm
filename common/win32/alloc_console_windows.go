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
	"io"
	"os"
	"syscall"

	"github.com/hkmh223/pd2mm/common/ansi"
)

// AllocConsole allocates a new console for the current process.
func AllocConsole() (aIn, aOut, aErr io.Writer, e error) { //nolint:nonamedreturns // reason: differentiate between writers.
	kernal32 := syscall.NewLazyDLL("kernel32.dll")
	allocConsole := kernal32.NewProc("AllocConsole")

	r0, _, err0 := syscall.SyscallN(allocConsole.Addr(), 0, 0, 0, 0)
	if r0 == 0 {
		return nil, nil, nil, fmt.Errorf("could not allocate console: %w. check build flags", err0)
	}

	hin, err1 := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	hout, err2 := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	herr, err3 := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)

	if err1 != nil {
		return nil, nil, nil, err1
	}

	if err2 != nil {
		return nil, nil, nil, err2
	}

	if err3 != nil {
		return nil, nil, nil, err3
	}

	stdinfile := os.NewFile(uintptr(hin), "/dev/stdin")
	stdoutfile := os.NewFile(uintptr(hout), "/dev/stdout")
	stderrfile := os.NewFile(uintptr(herr), "/dev/stderr")

	stdinA := ansi.NewAnsiStdoutW(stdinfile)
	stdoutA := ansi.NewAnsiStdoutW(stdoutfile)
	stderrA := ansi.NewAnsiStdoutW(stderrfile)

	os.Stdin = stdinfile
	os.Stdout = stdoutfile
	os.Stderr = stderrfile

	return stdinA, stdoutA, stderrA, nil
}
