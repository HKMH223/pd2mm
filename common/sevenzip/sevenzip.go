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

const (
	Name        = "7z"
	LinuxName   = "7z"
	WindowsName = "7z.exe"
)

// Extract extracts the contents of a 7z archive to a directory.
func Extract(src, dest string, redirect bool, opts ...ExtractionOptions) (ErrorCode, error) {
	opt := assureExtractionOptions(opts...)

	if !process.Exists(Name) {
		return ProcessNotFound, ErrSevenZipNotFound
	}

	if err := process.RunProcess(Name, opt.HideWindow, opt.Relative, redirect, "x", src, "-o"+dest+"/*"); err != nil {
		return CouldNotExtract, err
	}

	return NoError, nil
}

// ExtractWithBin extracts the contents of a 7z archive to a directory using a custom binary.
func ExtractWithBin(src, dest, bin string, redirect bool, opts ...ExtractionOptions) (ErrorCode, error) {
	opt := assureExtractionOptions(opts...)

	if !filesystem.Exists(bin) {
		return ProcessNotFound, ErrSevenZipNotFound
	}

	if err := process.RunProcess(bin, opt.HideWindow, opt.Relative, redirect, "x", src, "-o"+dest+"/*"); err != nil {
		return CouldNotCompress, err
	}

	return NoError, nil
}

// Compress compresses a directory to a 7z archive.
func Compress(src, dest string, redirect bool, opts ...CompressionOptions) (ErrorCode, error) {
	opt := assureCompressionOptions(opts...)

	if !process.Exists(Name) {
		return ProcessNotFound, ErrSevenZipNotFound
	}

	//nolint:lll // reason: calling RunProcess means that we can't pass CompressOptions directly.
	if err := process.RunProcess(Name, true, false, redirect, "a", "-t"+opt.FormatFormat, dest, src+"/*", opt.Level, opt.Method, opt.DictionarySize, opt.FastBytes, opt.SolidBlockSize, opt.Multithreading, opt.Memory); err != nil {
		return CouldNotCompress, err
	}

	return NoError, nil
}

// CompressWithBin compresses a directory to a 7z archive.
func CompressWithBin(src, dest, bin string, redirectStd bool, opts ...CompressionOptions) (ErrorCode, error) {
	opt := assureCompressionOptions(opts...)

	if !filesystem.Exists(bin) {
		return ProcessNotFound, ErrSevenZipNotFound
	}

	//nolint:lll // reason: calling RunProcess means that we can't pass CompressOptions directly.
	if err := process.RunProcess(bin, true, true, redirectStd, "a", "-t"+opt.FormatFormat, dest, src+"/*", opt.Level, opt.Method, opt.DictionarySize, opt.FastBytes, opt.SolidBlockSize, opt.Multithreading, opt.Memory); err != nil {
		return CouldNotCompress, err
	}

	return NoError, nil
}
