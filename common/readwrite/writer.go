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

import (
	"encoding/binary"
	"io"
	"os"
)

type Writer struct {
	file *os.File
}

type FileEntry struct {
	FileName      string
	FileNameLower uint32
	FileNameUpper uint32
	Offset        uint64
	UncompSize    uint64
}

type DataEntry struct {
	Hash     uint32
	FileName string
}

// FindByHash returns the first entry in data with a matching hash.
func FindByHash(data []DataEntry, hash uint32) *DataEntry {
	for _, entry := range data {
		if entry.Hash == hash {
			return &entry
		}
	}

	return nil
}

// FindByFileName returns the first entry in data with a matching file name.
func FindByFileName(data []DataEntry, name string) *DataEntry {
	for _, entry := range data {
		if entry.FileName == name {
			return &entry
		}
	}

	return nil
}

// NewWriter creates a new writer for the given file name.
func NewWriter(name string, appendMode bool) (*Writer, error) {
	var file *os.File

	var err error

	if appendMode {
		file, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	} else {
		file, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	}

	if err != nil {
		return nil, err
	}

	return &Writer{file}, nil
}

// WriteUInt32 writes a 32-bit unsigned integer to the file.
func (w *Writer) WriteUInt32(value uint32) error {
	return binary.Write(w.file, binary.LittleEndian, value)
}

// WriteUInt64 writes a 64-bit unsigned integer to the file.
func (w *Writer) WriteUInt64(value uint64) error {
	return binary.Write(w.file, binary.LittleEndian, value)
}

// Write writes data to the writer.
func (w *Writer) Write(data []byte) (int, error) {
	return w.file.Write(data)
}

// WriteChar writes a single character to the writer.
func (w *Writer) WriteChar(data string) (int, error) {
	return w.file.WriteString(data)
}

// Seek sets the position of the writer.
func (w *Writer) Seek(position int64, whence int) (int64, error) {
	return w.file.Seek(position, whence)
}

// Set the position of the reader to the beginning of the file.
func (w *Writer) SeekFromBeginning(position int64) (int64, error) {
	return w.file.Seek(position, io.SeekStart)
}

// Set the position of the reader to the end of the file.
func (w *Writer) SeekFromEnd(position int64) (int64, error) {
	return w.file.Seek(position, io.SeekEnd)
}

// Set the position of the reader to a specific position in the file.
func (w *Writer) SeekFromCurrent(position int64) (int64, error) {
	return w.file.Seek(position, io.SeekCurrent)
}

// Get the current position of the writer.
func (w *Writer) Position() (int64, error) {
	return w.file.Seek(0, io.SeekCurrent)
}

// Get the size of the file.
func (w *Writer) Size() (int64, error) {
	cur, err := w.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	defer func() {
		if _, err := w.file.Seek(cur, io.SeekStart); err != nil {
			panic(err)
		}
	}()

	size, err := w.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	return size, nil
}

// Close the file.
func (w *Writer) Close() error {
	return w.file.Close()
}
