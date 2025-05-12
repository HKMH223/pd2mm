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
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hkmh223/pd2mm/common/logger"
)

func DefaultUnzipMessenger() Messenger {
	return Messenger{
		AddedFile: func(path string) {
			logger.SharedLogger.Infof("Unzipping file: %s", path)
		},
	}
}

func Unzip(src, dst string) error {
	return UnzipByPrefixWithMessenger(src, dst, "", DefaultUnzipMessenger())
}

func UnzipByPrefixWithMessenger(src, dst, prefix string, msg Messenger) error {
	step := 1024
	read, err := zip.OpenReader(src)
	if err != nil { //nolint:wsl // gofumpt conflict
		return err
	}

	for _, file := range read.File {
		tmp, err := file.Open()
		if err != nil {
			return err
		}
		defer tmp.Close()

		if prefix != "" && !strings.HasPrefix(file.Name, prefix) {
			continue
		}

		name, err := url.QueryUnescape(maybeTrimPrefix(file.Name, prefix))
		if err != nil {
			return err
		}

		path := filepath.Join(dst, name)
		msg.AddedFile(path)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				continue
			}
		} else if err := stepCopy(path, tmp, int64(step)); err != nil {
			return err
		}
	}

	read.Close()

	return nil
}

func maybeTrimPrefix(trim, prefix string) string {
	if prefix != "" {
		return strings.TrimPrefix(trim, prefix)
	}

	return trim
}

func stepCopy(src string, dst io.Reader, step int64) error {
	if err := os.MkdirAll(filepath.Dir(src), os.ModePerm); err != nil {
		return err
	}

	path, err := os.Create(src)
	if err != nil {
		return err
	}

	for {
		_, err := io.CopyN(path, dst, step)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}
	}

	path.Close()

	return nil
}
