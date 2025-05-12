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
	"crypto/md5" //nolint:gosec // allowed
	"os"
	"path/filepath"
)

type DiffData struct {
	Hashes DiffHashData
	Local  DiffLocalData
}

type DiffHashData struct {
	File  string
	PathA string
	PathB string
	HashA string
	HashB string
}

type DiffLocalData struct {
	Path    string
	ExistsA string
	ExistsB string
}

func HashDirectory(dir string) (map[string]string, error) {
	hashes := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			rel, _ := filepath.Rel(dir, path)

			hash, err := NewHash(path, md5.New()) //nolint:gosec // allowed
			if err != nil {
				return err
			}

			hashes[rel] = hash
		}

		return nil
	})

	return hashes, err
}

func DiffDirectory(hashesA, hashesB map[string]string, dirA, dirB string) []DiffData {
	var diff []DiffData

	for pathA, hashA := range hashesA {
		if hashB, exists := hashesB[pathA]; !exists {
			diff = append(diff, DiffData{
				DiffHashData{}, //nolint:exhaustruct // allowed
				DiffLocalData{
					Path:    pathA,
					ExistsA: dirA,
					ExistsB: dirB,
				},
			})
		} else if hashA != hashB {
			diff = append(diff, DiffData{
				DiffHashData{
					File:  pathA,
					PathA: dirA,
					PathB: dirB,
					HashA: hashA,
					HashB: hashB,
				}, DiffLocalData{}, //nolint:exhaustruct // allowed
			})
		}
	}

	for pathB := range hashesB {
		if _, exists := hashesA[pathB]; !exists {
			diff = append(diff, DiffData{
				DiffHashData{}, //nolint:exhaustruct // allowed
				DiffLocalData{
					Path:    pathB,
					ExistsA: dirB,
					ExistsB: dirA,
				},
			})
		}
	}

	return diff
}
