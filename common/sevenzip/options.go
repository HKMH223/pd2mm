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

package sevenzip

type Options struct {
	FormatFormat   string
	Level          string
	Method         string
	DictionarySize string
	FastBytes      string
	SolidBlockSize string
	Multithreading string
	Memory         string
}

func getDefaultOptions() Options {
	return Options{
		FormatFormat:   "7z",
		Level:          "-mx9",
		Method:         "-m0=lzma2",
		DictionarySize: "-md=64m",
		FastBytes:      "-mfb=64",
		SolidBlockSize: "-ms=4g",
		Multithreading: "-mmt=2",
		Memory:         "-mmemuse=26g",
	}
}

func assureOptions(opts ...Options) Options {
	defopts := getDefaultOptions()

	if len(opts) == 0 {
		return defopts
	}

	return opts[0]
}
