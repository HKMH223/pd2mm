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
	"unicode/utf16"
	"unsafe"
)

func StringToBytes(str string) []byte {
	tmp := []byte(str)
	tmp = append(tmp, bytes.Repeat([]byte{0}, 2-(len(tmp)%2))...) //nolint:mnd // allowed

	return tmp
}

func GetStringFromBytes(data []byte, start, end int) string {
	var cid string

	if end > len(data) {
		end = len(data)
	}

	raw := data[start:end]
	u16 := ((*[1 << 30]uint16)(unsafe.Pointer(&raw[0])))[:len(raw)/2]

	ind := -1

	for i, c := range u16 {
		if c == 0 {
			ind = i
			break
		}
	}

	if ind != -1 {
		cid = string(utf16.Decode(u16[:ind]))
	} else {
		cid = string(utf16.Decode(u16))
	}

	return cid
}
