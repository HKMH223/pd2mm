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
	"errors"
	"flag"
	"fmt"
	"slices"
	"strings"
)

const BinSize int = 2

var ErrNoFunctionName = errors.New("no function name")

// IsFlagPassed returns true if the flag is passed.
func IsFlagPassed(name string) bool {
	var found bool

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}

// NewCommand returns a new command from the provided arguments.
func NewCommand(args []string, name string, count int) ([]string, error) {
	var names []string

	if slices.Contains(args, name) {
		ind := slices.Index(args, name)

		if ind != -1 {
			if len(args) >= count+BinSize {
				names = append(names, args[ind+1:]...)
				return names, nil
			}

			return nil, errExpectedArgs(count, len(args)-BinSize)
		}

		return nil, fmt.Errorf("%s not found", name) //nolint:err113 // reason: name must be known.
	}

	return nil, ErrNoFunctionName
}

// SplitArguments splits a string into arguments.
func SplitArguments(str string) []string {
	var parts []string

	var part strings.Builder

	var quote bool

	for _, char := range str {
		switch {
		case char == '"':
			quote = !quote
		case char == ' ' && !quote:
			parts = append(parts, part.String())
			part.Reset()
		default:
			part.WriteRune(char)
		}
	}

	parts = append(parts, part.String())

	return parts
}

// CheckArgumentCount checks if the number of arguments is correct.
func CheckArgumentCount(args []string, expected int) error {
	if len(args) != expected {
		return errExpectedArgs(expected, len(args))
	}

	return nil
}

// errexpectedArgs returns an error indicating that the number of arguments is incorrect.
func errExpectedArgs(expected, got int) error {
	return fmt.Errorf("expected: %d arguments but got %d", expected, got) //nolint:err113 // reason: values must be known.
}
