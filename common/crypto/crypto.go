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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"os"
	"strings"
)

const (
	CipherTag    = "<CipherKey>"
	Base64_16Len = 24
)

var DlfKey = []byte{ //nolint:gochecknoglobals // reason: DLF key is constant.
	65, 50, 114, 45, 208, 130, 239, 176, 220, 100, 87, 197, 118, 104, 202, 9,
}

var IV = make([]byte, 16) //nolint:gochecknoglobals,mnd // reason: IV is constant.

var (
	errDlfFileNotFound   = errors.New("error: DLF file not found")
	errCipherTagNotFound = errors.New("error: Cipher tag not found")
	errInvalidBase64Key  = errors.New("error: invalid base64 key")
	errInvalidIVSize     = errors.New("error: invalid IV size")
	errInvalidBufferSize = errors.New("error: invalid buffer size")
)

// NewHash creates a hash of the file at path.
func NewHash(path string, hash hash.Hash) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// AESEncrypt encrypts the text using the key and returns the encrypted data.
func AESEncrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	data := make([]byte, aes.BlockSize+len(text))
	iv := data[:aes.BlockSize] //nolint:varnamelen // reason: variable name used by cipher.NewCBCEncrypter().

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(data[aes.BlockSize:], text)

	return data, nil
}

// AESDecrypt decrypts the text using the key and returns the decrypted data.
func AESDecrypt(key, iv, buf []byte) ([]byte, error) { //nolint:varnamelen // reason: variable name used by cipher.NewCBCEncrypter().
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != aes.BlockSize {
		return nil, err
	}

	if len(buf) < aes.BlockSize {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)

	return buf, nil
}

// AESDecryptBase64 decrypts the text using the key and returns the decrypted data.
func AESDecryptBase64(kb64 string, iv, buf []byte) error { //nolint:varnamelen // reason: variable name used by cipher.NewCBCEncrypter().
	key, err := base64.StdEncoding.DecodeString(kb64)
	if err != nil {
		return errInvalidBase64Key
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	if len(iv) != aes.BlockSize {
		return errInvalidIVSize
	}

	if len(buf) < aes.BlockSize {
		return errInvalidBufferSize
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)

	return nil
}

// AESDecrypt decrypts the text using the key and returns the decrypted data.
func DecryptDLF(data []byte) ([]byte, error) {
	decrypt, err := AESDecrypt(DlfKey, data[0x41:], []byte{0})
	if err != nil {
		return nil, err
	}

	return decrypt, nil
}

// GetDLFAuto returns the decrypted data for the given CID.
func GetDLFAuto(cid string) ([]byte, error) {
	paths := []string{
		cid + ".dlf",
		cid + "_cached.dlf",
	}

	for _, path := range paths {
		data, err := readfile(path)
		if err == nil {
			return DecryptDLF(data)
		}
	}

	return nil, errDlfFileNotFound
}

// DecodeCipherTag decodes the cipher tag from the given data.
func DecodeCipherTag(dlf []byte) ([]byte, error) {
	data := string(dlf)
	pos := strings.Index(data, CipherTag)

	if pos == -1 {
		return nil, errCipherTagNotFound
	}

	pos += len(CipherTag)
	b64 := data[pos : pos+Base64_16Len]

	decode, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	if len(decode) > 16 { //nolint:mnd // reason: size of the cipher tag.
		decode = decode[:16]
	}

	return decode, nil
}

// GetOoaHash returns the OOA hash from the given data.
func GetOoaHash(data []byte) []byte {
	if len(data) < 0x3E { //nolint:mnd // reason: size of the OOA hash.
		return nil
	}

	return data[0x2A:0x3E]
}

// ReadFile reads the given file and returns its contents.
func readfile(name string) ([]byte, error) {
	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}
