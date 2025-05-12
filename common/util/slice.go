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

import "slices"

func Matches(parts []string, expected []string) int {
	expect := make(map[string]struct{}, len(expected))
	for _, e := range expected {
		expect[e] = struct{}{}
	}

	count := 0

	for _, part := range parts {
		if _, found := expect[part]; found {
			count++
		}
	}

	return count
}

func ReplaceSubslice[T comparable](slice, sliceA, sliceB []T) []T {
	for i := 0; i <= len(slice)-len(sliceA); i++ {
		if slices.Equal(slice[i:i+len(sliceA)], sliceA) {
			return slices.Concat(slice[:i], append(sliceB, slice[i+len(sliceA):]...))
		}
	}

	return slice
}

func ContainsSubslice[T comparable](sliceA, sliceB []T) bool {
	if len(sliceB) == 0 || len(sliceB) > len(sliceA) {
		return false
	}

	for i := 0; i <= len(sliceA)-len(sliceB); i++ {
		match := true

		for j := range sliceB {
			if sliceA[i+j] != sliceB[j] {
				match = false
				break
			}
		}

		if match {
			return true
		}
	}

	return false
}

func MoveEntry[T comparable](slice []T, entry T, index int) []T {
	cur := -1

	for i, s := range slice {
		if s == entry {
			cur = i
			break
		}
	}

	if cur == -1 {
		return slice
	}

	slice = slices.Delete(slice, cur, cur+1)

	if index >= len(slice) {
		slice = append(slice, entry)
	} else {
		slice = append(slice[:index], append([]T{entry}, slice[index:]...)...)
	}

	return slice
}
