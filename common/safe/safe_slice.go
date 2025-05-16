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

package safe

import (
	"fmt"
	"slices"

	"github.com/hkmh223/pd2mm/common/logger"
)

// HasIndex returns the index of a value in a slice.
// If the value is not found in the slice, an error will be logged and the program will exit.
func HasIndex[S ~[]E, E comparable](s S, v E) int {
	index := slices.Index(s, v)
	if index == -1 {
		logger.SharedLogger.Fatalf("%v expected %v but it was not found (index -1)", s, v)
	}

	return index
}

// Slice bounds checks to check if a slice is within bounds.
// If the slice is out of bounds, an error will be logged and the program will exit.
// Otherwise, the slice will be returned.
func Slice[T any](parts []T, index int) T { //nolint:ireturn // reason: T should not have constraints.
	return RangeWithCaller(parts, index, index+1, defaultCaller)[0] // There will only ever be one element in the returned slice.
} //nolint:ireturn // reason: T should not have constraints.

// Range bounds checks to check if a slice is within bounds.
// If the slice is out of bounds, an error will be logged and the program will exit.
// Otherwise, the slice will be returned.
func Range[T any](parts []T, start, end int) []T {
	return RangeWithCaller(parts, start, end, defaultCaller)
}

// SliceWithCaller bounds checks to check if a slice is within bounds.
// If the slice is out of bounds, an error will be logged and the program will exit.
// Otherwise, the slice will be returned.
//

func SliceWithCaller[T any](parts []T, index int, caller func(string)) T { //nolint:ireturn // reason: T should not have constraints.
	return RangeWithCaller(parts, index, index+1, caller)[0] // There will only ever be one element in the returned slice.
} //nolint:ireturn // reason: T should not have constraints.

// RangeWithCaller bounds checks to check if a slice is within bounds.
// If the slice is out of bounds, an error will be logged and the program will exit.
// Otherwise, the slice will be returned.
func RangeWithCaller[T any](parts []T, start, end int, caller func(string)) []T {
	if start < 0 || end > len(parts) || start > end {
		caller(fmt.Sprintf("invalid slice bounds: start=%d, end=%d, len=%d", start, end, len(parts)))
		return nil
	}

	return parts[start:end]
}

func defaultCaller(msg string) {
	logger.SharedLogger.Fatal(msg)
}
