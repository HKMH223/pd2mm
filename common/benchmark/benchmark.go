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

package benchmark

import "time"

// TimerWithResult with result.
//
//nolint:ireturn // allowed
func TimerWithResult[T any](fn func() (T, error), methodName string, caller func(string, string)) (T, error) {
	start := time.Now()
	result, err := fn()
	elapsed := time.Since(start)
	caller(methodName, elapsed.String())

	return result, err
}

// TimeAsync without result.
func Timer(fn func() error, methodName string, caller func(string, string)) error {
	start := time.Now()
	err := fn()
	elapsed := time.Since(start)
	caller(methodName, elapsed.String())

	return err
}
