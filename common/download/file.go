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

package download

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hkmh223/pd2mm/common/logger"
)

var (
	errDownloadURLEmpty  = errors.New("download url is empty")
	errDownloadPathEmpty = errors.New("download path is empty")
	errDownloadNameEmpty = errors.New("download name is empty")
	errFileHashNoMatch   = errors.New("file hash does not match")
)

type Messenger struct {
	StartDownload func(string)
}

func DefaultDownloadMessenger() Messenger {
	return Messenger{
		StartDownload: func(name string) {
			logger.SharedLogger.Infof("%s ... DOWNLOADING", name)
		},
	}
}

func DefaultHashValidator(path, hash, name string) error {
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err == nil {
			sha := sha256.New()
			sha.Write(data)
			sum := hex.EncodeToString(sha.Sum(nil))

			if strings.ToLower(hash) == sum {
				logger.SharedLogger.Infof("%s ... OK", name)
				return nil
			}
		}
	}

	return errFileHashNoMatch
}

func File(url, name, path string) error {
	return FileWithContext(context.TODO(), DefaultDownloadMessenger(), url, "", name, path, nil)
}

func FileValidated(url, hash, name, path string) error {
	return FileWithContext(context.TODO(), DefaultDownloadMessenger(), url, hash, name, path, DefaultHashValidator)
}

func FileWithBytes(url, name, path string) ([]byte, error) {
	if err := FileWithContext(context.TODO(), DefaultDownloadMessenger(), url, "", name, path, nil); err != nil {
		return nil, err
	}

	return read(path, name)
}

func FileWithBytesValidated(url, hash, name, path string) ([]byte, error) {
	if err := FileWithContext(context.TODO(), DefaultDownloadMessenger(), url, hash, name, path, DefaultHashValidator); err != nil {
		return nil, err
	}

	return read(path, name)
}

//nolint:lll // allowed
func FileWithContextAndBytes(ctx context.Context, state Messenger, url, hash, name, path string, validator func(string, string, string) error) ([]byte, error) {
	if err := FileWithContext(ctx, state, url, hash, name, path, validator); err != nil {
		return nil, err
	}

	return read(path, name)
}

//nolint:lll // allowed
func FileWithContext(ctx context.Context, state Messenger, url, hash, name, path string, validator func(string, string, string) error) error {
	if err := validateDownloadParams(url, path, name); err != nil {
		return err
	}

	fpath := filepath.Join(path, name)
	if err := os.MkdirAll(filepath.Dir(fpath), 0o700); err != nil {
		return err
	}

	if validator != nil {
		if err := validator(fpath, hash, name); err == nil {
			return nil
		}
	}

	state.StartDownload(name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	if validator != nil {
		return write(res, file, hash, name, false)
	}

	return write(res, file, hash, name, true)
}

func validateDownloadParams(url, path, name string) error {
	if url == "" {
		return errDownloadURLEmpty
	}

	if path == "" {
		return errDownloadPathEmpty
	}

	if name == "" {
		return errDownloadNameEmpty
	}

	return nil
}

func read(path, name string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(path, name))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func write(resp *http.Response, flags *os.File, hash, name string, skip bool) error {
	sha := sha256.New()
	buf := make([]byte, 1<<20) //nolint:mnd // 1 megabyte buffer

	for {
		index, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if index == 0 {
			break
		}

		if _, err := flags.Write(buf[:index]); err != nil {
			return err
		}

		if _, err := sha.Write(buf[:index]); err != nil {
			return err
		}
	}

	if skip {
		return nil
	}

	sum := hex.EncodeToString(sha.Sum(nil))
	if strings.ToLower(hash) != sum {
		return fmt.Errorf("hash mismatch for %s", name) //nolint:err113 // required name prevents static error
	}

	return nil
}
