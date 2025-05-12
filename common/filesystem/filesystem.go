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

package filesystem

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/otiai10/copy"
)

var (
	errNoNameFound = errors.New("name could not be found in file path")
	errNoPathFound = errors.New("file path could not be found in file name")
	errFileExists  = errors.New("file exists in destination path")
)

var ReservedHostnames = []string{ //nolint:gochecknoglobals // allowed
	"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
	"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	"PRN", "AUX", "NUL",
}

func Combine(pathA string, pathB ...string) string {
	path := append([]string{pathA}, pathB...)
	return filepath.Join(path...)
}

func FromCwd(pathA ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := append([]string{wd}, pathA...)

	return filepath.Join(path...), nil
}

func GetDirectoryName(name string) string {
	return filepath.Dir(name)
}

func GetFileName(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func GetFileExtension(name string) string {
	return filepath.Ext(name)
}

func GetRelativePath(paths ...string) string {
	result := "./" + paths[0]

	for _, dir := range paths[1:] {
		result = path.Join(result, dir)
	}

	return result
}

func TrimPath(path string) string {
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, ".\\") {
		return path[2:]
	} else if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
		return path[1:]
	}

	return path
}

func Copy(pathA, pathB string, opts ...copy.Options) error {
	if err := copy.Copy(pathA, pathB, opts...); err != nil {
		return err
	}

	return nil
}

func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)

	return err
}

func CopyAndRename(files []string, oldPath, newPath, oldName, newName string) error {
	found := false

	for _, file := range files {
		if strings.Contains(file, oldName) {
			found = true
			break
		}
	}

	if !found {
		return errNoNameFound
	}

	for _, file := range files {
		newName := strings.ReplaceAll(file, oldName, newName)

		if !strings.Contains(newName, TrimPath(oldPath)) {
			return errNoPathFound
		}

		newFilePath := strings.ReplaceAll(newName, TrimPath(oldPath), newPath)

		if Exists(newFilePath) {
			return errFileExists
		}

		if err := Copy(file, newFilePath); err != nil {
			return err
		}
	}

	return nil
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func ReadFile(name string) ([]byte, error) {
	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func ReadAllLines(file *os.File) ([]string, error) {
	return Scan(bufio.NewScanner(file))
}

func ReadAllStringLines(str string) ([]string, error) {
	return Scan(bufio.NewScanner(strings.NewReader(str)))
}

func WriteFile(name string, data []byte, perm fs.FileMode) error {
	err := os.WriteFile(name, data, perm)
	if err != nil {
		return err
	}

	return nil
}

func WriteToFile(file *os.File, entries []string) error {
	for _, entry := range entries {
		if _, err := file.WriteString(entry); err != nil {
			return err
		}
	}

	return nil
}

func OverwriteFile(file *os.File) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}

func Scan(scanner *bufio.Scanner) ([]string, error) {
	lines := []string{}

	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}

		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func DeleteDirectory(name string) error {
	err := os.RemoveAll(name)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEmptyDirectories(dir string) error {
	directories := []string{}

	err := filepath.WalkDir(dir, func(path string, directory os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if directory.IsDir() && path != dir {
			directories = append(directories, path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	for i := len(directories) - 1; i >= 0; i-- {
		dir := directories[i]

		empty, err := IsEmpty(dir)
		if err != nil {
			return err
		}

		if empty {
			err = os.Remove(dir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SortFileNames(paths []string) []string {
	sort.Slice(paths, func(i, j int) bool {
		parentA := filepath.Dir(paths[i])
		parentB := filepath.Dir(paths[j])

		if parentA == parentB {
			return filepath.Base(paths[i]) < filepath.Base(paths[j])
		}

		return parentA < parentB
	})

	return paths
}

func GetFiles(path string) []string {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return []string{}
	}

	return SortFileNames(files)
}

func GetDirectories(path string) []string {
	var directories []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			directories = append(directories, path)
		}

		return nil
	})
	if err != nil {
		return []string{}
	}

	return SortFileNames(directories)
}

func GetTopDirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var directories []string

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}

	return directories, nil
}

func IsEmpty(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.Readdir(1)

	if err == nil {
		return false, nil
	}

	if errors.Is(err, os.ErrNotExist) || err.Error() == "EOF" {
		return true, nil
	}

	return false, err
}

func BytesToMap(data []byte) (map[string]any, error) {
	var b map[string]any
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}

	return b, nil
}

func FilenameToMap(initial, name string) (map[string]any, error) {
	data, err := os.ReadFile(initial + name)
	if err != nil {
		return nil, err
	}

	var b map[string]any
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}

	return b, nil
}

func FilenameToBytes(initial, name string) ([]byte, error) {
	data, err := os.ReadFile(initial + name)
	if err != nil {
		return nil, err
	}

	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	return data, nil
}

func IsValidHostname(hostname string) bool {
	if len(hostname) < 1 || len(hostname) > 15 {
		return false
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

	if !regex.MatchString(hostname) {
		return false
	}

	for _, reserved := range ReservedHostnames {
		if strings.EqualFold(hostname, reserved) {
			return false
		}
	}

	return true
}
