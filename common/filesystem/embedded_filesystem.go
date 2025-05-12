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
	"embed"
	"encoding/json"
	"fmt"
)

type EmbeddedFileSystem struct {
	Initial string
	FS      embed.FS
}

// BytesToMap converts a byte array into a map[string]any.
func (e EmbeddedFileSystem) BytesToMap(data []byte) (map[string]any, error) {
	var b map[string]any
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}

	return b, nil
}

// FilenameToMap converts a file in the embedded filesystem into a map[string]any.
func (e EmbeddedFileSystem) FilenameToMap(name string) (map[string]any, error) {
	data, err := e.FS.ReadFile(e.Initial + name)
	if err != nil {
		return nil, err
	}

	var b map[string]any
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}

	return b, nil
}

// FilenameToBytes converts a file in the embedded filesystem into an array of bytes.
func (e EmbeddedFileSystem) FilenameToBytes(name string) ([]byte, error) {
	data, err := e.FS.ReadFile(e.Initial + name)
	if err != nil {
		return nil, err
	}

	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	return data, nil
}
