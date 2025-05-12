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
	"fmt"
	"regexp"
	"strings"
)

func DrawWatermark(text []string, draw func(string)) {
	result := []string{}

	longest := 0

	for _, txt := range text {
		length := textLength(txt)
		if length > longest {
			longest = length
		}
	}

	line := strings.Repeat("-", longest)
	result = append(result, fmt.Sprintf("┌─%s─┐", line))

	for _, txt := range text {
		spaceSize := longest - textLength(txt)
		spaceText := txt + strings.Repeat(" ", spaceSize)
		result = append(result, fmt.Sprintf("│ %s │", spaceText))
	}

	result = append(result, fmt.Sprintf("└─%s─┘", line))

	for _, txt := range result {
		draw(txt)
	}
}

func textLength(s string) int {
	re := regexp.MustCompile(`[\p{Han}\p{Katakana}\p{Hiragana}\p{Hangul}]`)
	result := re.ReplaceAllString(s, "ab")

	return len(result)
}
