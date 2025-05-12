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

package util

import (
	"bytes"
	"fmt"
)

func HexStringToBytes(hex string) ([]byte, error) {
	var data []byte

	for i := 0; i < len(hex); i += 2 {
		var bytes byte

		_, err := fmt.Sscanf(hex[i:i+2], "%02X", &bytes)
		if err != nil {
			return nil, err
		}

		data = append(data, bytes)
	}

	return data, nil
}

func FindAllByteOccurrences(data []byte, pattern []byte) []int {
	var ind []int

	for i := range data {
		if bytes.HasPrefix(data[i:], pattern) {
			ind = append(ind, i)
		}
	}

	return ind
}

func ReplaceByteOccurrences(original []byte, expected []byte, replacement []byte, occurrence int) []byte {
	var result []byte

	remaining := original
	count := 0

	for {
		index := bytes.Index(remaining, expected)
		if index == -1 {
			result = append(result, remaining...)
			break
		}

		result = append(result, remaining[:index]...)

		count++

		if occurrence == 0 || count == occurrence {
			txt := min(len(replacement), len(expected))

			result = append(result, replacement[:txt]...)
		} else {
			result = append(result, remaining[index:index+len(expected)]...)
		}

		remaining = remaining[index+len(expected):]
	}

	return result
}
