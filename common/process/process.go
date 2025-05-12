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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func DoesFileExist(name string) bool {
	if _, err := exec.LookPath(name); err != nil {
		return false
	}

	return true
}

func RunFile(name string, hide, rel, redirect bool, arg ...string) error {
	path := name

	if rel {
		cwd, err := os.Executable()
		if err != nil {
			return err
		}

		path = filepath.Join(filepath.Dir(cwd), name)
	}

	cmd := exec.Command(path, arg...)

	if redirect {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if runtime.GOOS == "windows" {
		setHideWindowAttr(cmd, hide)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
