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

package zip

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
)

type Messenger struct {
	AddedFile func(string)
}

// DefaultZipMessenger returns a default messenger for the Zip function.
func DefaultZipMessenger() Messenger {
	return Messenger{
		AddedFile: func(path string) {
			logger.SharedLogger.Infof("Adding file to zip: %s", path)
		},
	}
}

// Zip creates a zip file from the given source directory and destination path.
func Zip(src, dest string) error {
	return WithMessenger(src, dest, DefaultZipMessenger())
}

// WithContext is a convenience function that zips the file.
func WithMessenger(src, dest string, msg Messenger) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	write := zip.NewWriter(file)
	defer write.Close()

	err = filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		create, err := write.Create(convertPath(path, src))
		if err != nil {
			return err
		}

		msg.AddedFile(path)

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		if _, err = io.Copy(create, file); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func convertPath(path, src string) string {
	path = trimSrcPrefix(path, src)
	path = filesystem.Normalize(path)

	return path
}

func trimSrcPrefix(path, src string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, src), string(filepath.Separator))
}
