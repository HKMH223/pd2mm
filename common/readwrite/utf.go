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

package readwrite

import "unicode/utf16"

func Utf8ToUtf16(utf8 string) []byte {
	bytes := []byte(utf8)
	u16Runes := utf16.Encode([]rune(string(bytes)))
	u16Bytes := make([]byte, len(u16Runes)*2) //nolint:mnd // testcase

	for i, r := range u16Runes {
		u16Bytes[i*2] = byte(r)
		u16Bytes[i*2+1] = byte(r >> 8) //nolint:mnd // testcase
	}

	return u16Bytes
}
