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

package main

import (
	"crypto/md5" //nolint:gosec // allowed
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New() //nolint:gosec // allowed

	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

//nolint:cyclop // allowed
func compareFolders(folderA, folderB string) error {
	checksumsA := make(map[string]string)
	checksumsB := make(map[string]string)

	err := filepath.Walk(folderA, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relativePath, _ := filepath.Rel(folderA, path)

			checksum, err := calculateChecksum(path)
			if err != nil {
				return err
			}

			checksumsA[relativePath] = checksum
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = filepath.Walk(folderB, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relativePath, _ := filepath.Rel(folderB, path)

			checksum, err := calculateChecksum(path)
			if err != nil {
				return err
			}

			checksumsB[relativePath] = checksum
		}

		return nil
	})
	if err != nil {
		return err
	}

	for pathA, checksumA := range checksumsA {
		checksumB, exists := checksumsB[pathA]
		if !exists {
			fmt.Printf("File %s exists in folder1 but not in folder2\n", pathA)
			continue
		}

		if checksumA != checksumB {
			fmt.Printf("Checksums for file %s do not match:\n", pathA)
			fmt.Printf("  Folder1: %s\n", checksumA)
			fmt.Printf("  Folder2: %s\n", checksumB)
		}
	}

	for pathB := range checksumsB {
		_, exists := checksumsA[pathB]
		if !exists {
			fmt.Printf("File %s exists in folder2 but not in folder1\n", pathB)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 { //nolint:mnd // allowed
		fmt.Println("Usage: diff <folder1> <folder2>")
		os.Exit(1)
	}

	folderA := os.Args[1]
	folderB := os.Args[2]

	err := compareFolders(folderA, folderB)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
