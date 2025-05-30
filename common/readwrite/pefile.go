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
	"debug/pe"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidOffsetOrByteRange = errors.New("invalid offset or byte range")
	ErrSectionHeaderIsSizeZero  = errors.New("section header size is 0")
	ErrSectionIsNil             = errors.New("section is nil")
	ErrNoBytes                  = errors.New("no bytes")
)

// COFFHeader
// 0x50, 0x45 = PE
// COFF_START_BYTES_LEN == len(COFFStartBytes).
var COFFStartBytes = []byte{0x50, 0x45, 0x00, 0x00} //nolint:gochecknoglobals // reason: COFFStartBytes is constant.

const (
	COFFStartBytesLen = 4
	COFFHeaderSize    = 20
)

// OptionalHeader64
// https://github.com/golang/go/blob/master/src/debug/pe/pe.go
// uint byte size of OptionalHeader64 without magic mumber(2 bytes) or data directory(128 bytes)
// OptionalHeader64 size is 240
// (110).
var OH64ByteSize = binary.Size(OptionalHeader64X110{}) //nolint:exhaustruct,gochecknoglobals // reason: OH64ByteSize is constant.

// DataDirectory
// 16 entries * 8 bytes / entry.
const (
	DataDirSize      = 128
	DataDirEntrySize = 8
)

// SectionHeader32
// https://github.com/golang/go/blob/master/src/debug/pe/section.go
// uint byte size of SectionHeader32 without name(8 bytes) or characteristics(4 bytes)
// (28).
var SH32ByteSize = binary.Size(SectionHeader32X28{}) //nolint:exhaustruct,gochecknoglobals // reason: SH32ByteSize is constant.

const (
	SH32EntrySize           = 64
	SH32NameSize            = 8
	SH32CharacteristicsSize = 4
)

// Data structure.
type Data struct {
	Bytes []byte
	PE    pe.File
}

// Section structure (.ooa).
type Section struct {
	ContentID   string
	OEP         uint64
	EncBlocks   []EncBlock
	ImageBase   uint64
	SizeOfImage uint32
	ImportDir   DataDir
	IATDir      DataDir
	RelocDir    DataDir
}

// Import structure (.ooa).
type Import struct {
	Characteristics uint32
	Timedatestamp   uint32
	ForwarderChain  uint32
	Name            uint32
	FThunk          uint32
}

// Thunk structure (.ooa).
type Thunk struct {
	Function uint32
	DataAddr uint32
}

// DataDir structure (.ooa).
type DataDir struct {
	VA   uint32
	Size uint32
}

// EncBlock structure (.ooa).
type EncBlock struct {
	VA          uint32
	RawSize     uint32
	VirtualSize uint32
	Unk         uint32
	CRC         uint32
	Unk2        uint32
	CRC2        uint32
	Pad         uint32
	FileOffset  uint32
	Pad2        uint64
	Pad3        uint32
}

type OptionalHeader64X110 struct {
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
}

type SectionHeader32X28 struct {
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLineNumbers uint32
	NumberOfRelocations  uint16
	NumberOfLineNumbers  uint16
}

var errRVALessThanOne = errors.New("rva is equal to '-1'")

// Open a file at the specified path and return Data.
func Open(path string) (*Data, error) {
	data := new(Data)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	pefile, err := pe.NewFile(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	data.Bytes = bytes
	data.PE = *pefile

	return data, nil
}

// WriteBytes writes the specified bytes to the file at the given offset.
func WriteBytes(data []byte, offset int, replace []byte) error {
	if offset < 0 || offset+len(replace) > len(data) {
		return ErrInvalidOffsetOrByteRange
	}

	copy(data[offset:], replace)

	return nil
}

// FindBytes searches for the specified bytes in the file.
func ReadCOFFHeaderOffset(data []byte) (int, error) {
	offset, err := FindBytes(data, COFFStartBytes)
	if err != nil {
		return -1, err
	}

	return offset, nil
}

// ReadDDBytes reads the data directory entry at the specified offset.
func ReadDDBytes(data []byte) ([]byte, error) {
	offset, err := ReadCOFFHeaderOffset(data)
	if err != nil {
		return nil, err
	}

	return data[offset+COFFStartBytesLen+COFFHeaderSize+OH64ByteSize : offset+COFFStartBytesLen+COFFHeaderSize+OH64ByteSize+DataDirSize], nil
}

// ReadDDEntryOffset reads the offset of the data directory entry at the specified address.
func ReadDDEntryOffset(data []byte, addr, size uint32) (int, error) {
	dir, err := ReadDDBytes(data)
	if err != nil {
		return -1, err
	}

	bytes := make([]byte, DataDirEntrySize)
	binary.LittleEndian.PutUint32(bytes[:4], addr)
	binary.LittleEndian.PutUint32(bytes[4:], size)
	rva, err := FindBytes(dir, bytes)

	if err != nil || rva == -1 {
		if err == nil {
			return -1, errRVALessThanOne
		}

		return -1, err
	}

	offset, err := ReadCOFFHeaderOffset(data)
	if err != nil {
		return -1, err
	}

	return offset + COFFStartBytesLen + COFFHeaderSize + OH64ByteSize + rva, nil
}

// ReadSHBytes reads the section header bytes at the specified offset.
func ReadSHSize(file pe.File) (int, error) {
	sections := len(file.Sections)
	size := sections * SH32EntrySize

	if size == 0 {
		return -1, ErrSectionHeaderIsSizeZero
	}

	return size, nil
}

// ReadSHEntry reads the section header entry at the specified offset.
func ReadSHBytes(data []byte, size int) ([]byte, error) {
	offset, err := ReadCOFFHeaderOffset(data)
	if err != nil {
		return nil, err
	}

	index := offset + COFFStartBytesLen + COFFHeaderSize + OH64ByteSize + DataDirSize

	return data[index : index+size], nil
}

// ReadSHEntryOffset reads the offset of the specified section header entry.
func ReadSHEntryOffset(data []byte, address int) (int, error) {
	offset, err := ReadCOFFHeaderOffset(data)
	if err != nil {
		return -1, err
	}

	return offset + COFFStartBytesLen + COFFHeaderSize + OH64ByteSize + DataDirSize + address, nil
}

// ReadSectionBytes reads the specified section bytes.
func ReadSectionBytes(file *Data, sectionVirtualAddress, sectionSize uint32) ([]byte, error) {
	var section *pe.Section

	for _, s := range file.PE.Sections {
		if sectionVirtualAddress >= s.VirtualAddress && sectionVirtualAddress < s.VirtualAddress+s.Size {
			section = s
			break
		}
	}

	if section == nil {
		return nil, ErrSectionIsNil
	}

	offset := sectionVirtualAddress - section.VirtualAddress + section.Offset
	bytes := file.Bytes[offset : offset+sectionSize]

	return bytes, nil
}

// ReadImport reads the import section.
func ReadImport(reader io.Reader) (Import, error) {
	var data Import
	err := binary.Read(reader, binary.LittleEndian, &data)

	return data, err
}

// ReadThunk reads the thunk section.
func ReadThunk(reader io.Reader) (Thunk, error) {
	var data Thunk
	err := binary.Read(reader, binary.LittleEndian, &data)

	return data, err
}

// ReadEncBlock reads the encryption block.
func ReadDataDir(reader io.Reader) (DataDir, error) {
	var data DataDir
	err := binary.Read(reader, binary.LittleEndian, &data)

	return data, err
}

// ReadEncBlock reads the encryption block.
func ReadEncBlock(reader io.Reader) (EncBlock, error) {
	var data EncBlock
	err := binary.Read(reader, binary.LittleEndian, &data)

	return data, err
}

// FindBytes finds the index of the first occurrence of dest in src.
func FindBytes(src, dest []byte) (int, error) {
	for i := range src[:len(src)-len(dest)+1] {
		if MatchBytes(src[i:i+len(dest)], dest) {
			return i, nil
		}
	}

	return -1, ErrNoBytes
}

// PadBytes pads the bytes to the specified size.
func PadBytes(data []byte, size int) []byte {
	if len(data) < size {
		paddingSize := size - len(data)
		padding := make([]byte, paddingSize)

		return append(data, padding...)
	}

	return data
}

// MatchBytes matches the bytes in src to the bytes in dest.
func MatchBytes(src, dest []byte) bool {
	for i := range dest {
		if src[i] != dest[i] {
			return false
		}
	}

	return true
}
