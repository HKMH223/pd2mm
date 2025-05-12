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

package crypto

import (
	"crypto/md5"  //nolint:gosec // allowed
	"crypto/sha1" //nolint:gosec // allowed
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"hash"
	"hash/crc32"
	"hash/crc64"

	"github.com/hkmh223/pd2mm/common/murmurhash3"
	"github.com/hkmh223/pd2mm/common/readwrite"
)

var errHashNotEqual = errors.New("file hash is not equal to specified hash")

func Validate(path, hash string, hashType hash.Hash) error {
	hashA, err := NewHash(path, hashType)
	if err != nil {
		return err
	}

	if hashA != hash {
		return errHashNotEqual
	}

	return nil
}

func NewMD5(path string) (string, error) {
	s, err := NewHash(path, md5.New()) //nolint:gosec // allowed
	if err != nil {
		return "", err
	}

	return s, nil
}

func NewSHA1(path string) (string, error) {
	s, err := NewHash(path, sha1.New()) //nolint:gosec // allowed
	if err != nil {
		return "", err
	}

	return s, nil
}

func NewSHA256(path string) (string, error) {
	s, err := NewHash(path, sha256.New())
	if err != nil {
		return "", err
	}

	return s, nil
}

func NewSHA512(path string) (string, error) {
	s, err := NewHash(path, sha512.New())
	if err != nil {
		return "", err
	}

	return s, nil
}

func NewCRC32(path string) (string, error) {
	s, err := NewHash(path, crc32.New(crc32.IEEETable))
	if err != nil {
		return "", err
	}

	return s, nil
}

func NewCRC64(path string) (string, error) {
	s, err := NewHash(path, crc64.New(crc64.MakeTable(crc32.IEEE)))
	if err != nil {
		return "", err
	}

	return s, nil
}

func Murmur3X64_128Hash(seed int, str string) uint64 {
	bytes := murmurhash3.NewX64_128(seed)
	bytes.Write(readwrite.Utf8ToUtf16(str))

	return binary.LittleEndian.Uint64(bytes.Sum(nil))
}

func Murmur3X86_32Hash(seed int, str string) uint32 {
	bytes := murmurhash3.NewX86_32(seed)
	bytes.Write(readwrite.Utf8ToUtf16(str))

	return binary.LittleEndian.Uint32(bytes.Sum(nil))
}

func Murmur3X86_128Hash(seed int, str string) uint32 {
	bytes := murmurhash3.NewX86_128(seed)
	bytes.Write(readwrite.Utf8ToUtf16(str))

	return binary.LittleEndian.Uint32(bytes.Sum(nil))
}
