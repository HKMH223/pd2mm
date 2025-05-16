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

package filesystem

import (
	"bytes"
	"strings"
	"sync"
)

type LineRingBuffer struct {
	mu       sync.Mutex
	lines    []string
	capacity int
	start    int
	count    int
	buf      bytes.Buffer
}

// NewLineRingBuffer creates a new LineRingBuffer with the given capacity.
func NewLineRingBuffer(capacity int) *LineRingBuffer {
	return &LineRingBuffer{ //nolint:exhaustruct // reason: set by functions
		lines:    make([]string, capacity),
		capacity: capacity,
	}
}

// Write appends the given data to the buffer.
func (r *LineRingBuffer) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	length := len(p)
	r.buf.Write(p)

	for {
		data := r.buf.Bytes()

		index := bytes.IndexByte(data, '\n')
		if index < 0 {
			break
		}

		line := string(data[:index])
		r.buf.Next(index + 1)

		r.addLine(line)
	}

	return length, nil
}

// addLine adds the given line to the buffer.
func (r *LineRingBuffer) addLine(line string) {
	if r.count < r.capacity {
		r.lines[(r.start+r.count)%r.capacity] = line
		r.count++

		return
	}

	r.lines[r.start] = line
	r.start = (r.start + 1) % r.capacity
}

// String returns the contents of the buffer as a string.
func (r *LineRingBuffer) String() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	var builder strings.Builder

	for i := range r.count {
		idx := (r.start + i) % r.capacity

		builder.WriteString(r.lines[idx])
		builder.WriteByte('\n')
	}

	return builder.String()
}

// Reset resets the buffer to its initial state.
func (r *LineRingBuffer) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lines = make([]string, r.capacity)
	r.start = 0
	r.count = 0
	r.buf.Reset()
}
