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

type Reader struct {
	file *os.File
}

// NewReader creates a new reader for the given file.
func NewReader(name string) (*Reader, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return &Reader{file}, nil
}

// IsValid returns true if the reader is valid.
func (r *Reader) IsValid() bool {
	return r.file != nil
}

// ReadUInt32 reads a 32-bit unsigned integer from the file.
func (r *Reader) ReadUInt32() (uint32, error) {
	var value uint32
	err := binary.Read(r.file, binary.LittleEndian, &value)

	return value, err
}

// ReadUInt64 reads a 64-bit unsigned integer from the file.
func (r *Reader) ReadUInt64() (uint64, error) {
	var value uint64
	err := binary.Read(r.file, binary.LittleEndian, &value)

	return value, err
}

// Read reads data from the reader.
func (r *Reader) Read(data []byte) (int, error) {
	return r.file.Read(data)
}

// ReadChar reads a single character from the reader.
func (r *Reader) ReadChar() (byte, error) {
	var value byte
	err := binary.Read(r.file, binary.LittleEndian, &value)

	return value, err
}

// Seek sets the position of the reader.
func (r *Reader) Seek(position int64, whence int) (int64, error) {
	return r.file.Seek(position, whence)
}

// Set the position of the reader to the beginning of the file.
func (r *Reader) SeekFromBeginning(position int64) (int64, error) {
	return r.file.Seek(position, io.SeekStart)
}

// Set the position of the reader to the end of the file.
func (r *Reader) SeekFromEnd(position int64) (int64, error) {
	return r.file.Seek(position, io.SeekEnd)
}

// Set the position of the reader to a specific position in the file.
func (r *Reader) SeekFromCurrent(position int64) (int64, error) {
	return r.file.Seek(position, io.SeekCurrent)
}

// Get the current position of the reader.
func (r *Reader) Position() (int64, error) {
	return r.file.Seek(0, io.SeekCurrent)
}

// Get the size of the file.
func (r *Reader) Size() (int64, error) {
	cur, err := r.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	defer func() {
		if _, err := r.file.Seek(cur, io.SeekStart); err != nil {
			panic(err)
		}
	}()

	size, err := r.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	return size, nil
}

// Close the reader.
func (r *Reader) Close() error {
	return r.file.Close()
}
