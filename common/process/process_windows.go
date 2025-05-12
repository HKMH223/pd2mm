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

package process

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

//nolint:gochecknoglobals // allowed
var (
	user32   = windows.NewLazyDLL("user32.dll")
	kernal32 = syscall.NewLazyDLL("kernel32.dll")
	psapi    = syscall.NewLazyDLL("psapi.dll")

	procGetWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess              = kernal32.NewProc("OpenProcess")
	procGetModuleFileNameEx      = psapi.NewProc("GetModuleFileNameExW")
	procGetWindowText            = user32.NewProc("GetWindowTextW")
	procGetWindowTextLength      = user32.NewProc("GetWindowTextLengthW")
)

var errFailedToGetProcessID = errors.New("failed to get process ID")

type (
	HANDLE uintptr
	HWND   HANDLE
)

const (
	ProcessQueryInformation = 0x0400
	ProcessVMRead           = 0x0010
)

// Get the window handle process id.
func GetWindowProcessID(hwnd HWND) (uint32, error) {
	var pid uint32

	tid, _, _ := procGetWindowThreadProcessID.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

	if tid == 0 {
		return 0, errFailedToGetProcessID
	}

	return pid, nil
}

// Get the window handle text length.
func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

// Get the window handle text.
func GetWindowText(hwnd HWND) (string, error) {
	txt := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, txt)

	_, _, err := procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(txt))
	if err != nil {
		return "", err
	}

	return syscall.UTF16ToString(buf), nil
}

// Get the window handle by function name.
func GetWindow(funcName string) uintptr {
	proc := user32.NewProc(funcName)
	hwnd, _, _ := proc.Call()

	return hwnd
}

// Get the process id executable name.
func GetExecutableName(pid uint32) (string, error) {
	hproc, _, err := procOpenProcess.Call(
		ProcessQueryInformation|ProcessVMRead,
		0,
		uintptr(pid),
	)

	if hproc == 0 {
		return "", fmt.Errorf("failed to open process: %w", err)
	}

	defer func() {
		if err := syscall.CloseHandle(syscall.Handle(hproc)); err != nil {
			log.Fatalf("Failed to close handle: %s", err)
		}
	}()

	buf := make([]uint16, 1024) //nolint:mnd // allowed

	ret, _, err := procGetModuleFileNameEx.Call(
		hproc,
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)

	if ret == 0 {
		return "", fmt.Errorf("failed to get module file name: %w", err)
	}

	return syscall.UTF16ToString(buf[:ret]), nil
}

// Get the executable name of the foreground window.
func GetForegroundWindowExecutableName() (string, error) {
	if hwnd := GetWindow("GetForegroundWindow"); hwnd != 0 {
		pid, err := GetWindowProcessID(HWND(hwnd))
		if err != nil {
			return "", fmt.Errorf("error getting process ID: %w", err)
		}

		exeName, err := GetExecutableName(pid)
		if err != nil {
			return "", fmt.Errorf("error getting executable name: %w", err)
		}

		return filepath.Base(exeName), nil
	}

	return "", nil
}

// Get the title of the foreground window.
func GetForegroundWindowTitle() (string, error) {
	if hwnd := GetWindow("GetForegroundWindow"); hwnd != 0 {
		title, err := GetWindowText(HWND(hwnd))
		if err != nil {
			return "", err
		}

		return title, nil
	}

	return "", nil
}

// Set the hide window attribute.
func setHideWindowAttr(cmd *exec.Cmd, hideWindow bool) {
	// cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: hideWindow} //nolint:exhaustruct // allowed
}
