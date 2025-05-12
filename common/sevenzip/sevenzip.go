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

import (
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/process"
)

// Extract extracts the contents of a 7z archive to a directory.
func Extract(src, dest string, redirect bool) (ErrorCode, error) {
	if !process.Exists("7z") {
		return ProcessNotFound, ErrSevenZipNotFound
	}

	if err := process.RunProcess("7z", true, false, redirect, "x", src, "-o"+dest+"/*"); err != nil {
		return CouldNotExtract, err
	}

	return NoError, nil
}

// ExtractFromBin extracts the contents of a 7z archive to a directory using a custom binary.
func ExtractFromBin(src, dest, bin string, redirect bool) (ErrorCode, error) {
	if !filesystem.Exists(bin) {
		return ProcessNotFound, ErrSevenZipNotFound
	}

	if err := process.RunProcess(bin, true, true, redirect, "x", src, "-o"+dest+"/*"); err != nil {
		return CouldNotCompress, err
	}

	return NoError, nil
}

// Compress compresses a directory to a 7z archive.
func Compress(src, dest string, redirect bool, opts ...Options) (ErrorCode, error) {
	opt := assureOptions(opts...)

	if !process.Exists("7z") {
		return ProcessNotFound, ErrSevenZipNotFound
	}

	//nolint:lll // allowed
	if err := process.RunProcess("7z", true, false, redirect, "a", "-t"+opt.FormatFormat, dest, src+"/*", opt.Level, opt.Method, opt.DictionarySize, opt.FastBytes, opt.SolidBlockSize, opt.Multithreading, opt.Memory); err != nil {
		return CouldNotCompress, err
	}

	return NoError, nil
}

// CompressFromBin compresses a directory to a 7z archive.
func CompressFromBin(src, dest, bin string, redirectStd bool, opts ...Options) (ErrorCode, error) {
	opt := assureOptions(opts...)

	if !filesystem.Exists(bin) {
		return ProcessNotFound, ErrSevenZipNotFound
	}

	//nolint:lll // allowed
	if err := process.RunProcess(bin, true, true, redirectStd, "a", "-t"+opt.FormatFormat, dest, src+"/*", opt.Level, opt.Method, opt.DictionarySize, opt.FastBytes, opt.SolidBlockSize, opt.Multithreading, opt.Memory); err != nil {
		return CouldNotCompress, err
	}

	return NoError, nil
}
