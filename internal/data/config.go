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

package data

import (
	"encoding/json"
	"os"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/tidwall/jsonc"
)

var FileTypes = []string{".jsonc", ".json"} //nolint:gochecknoglobals // reason: file types are needed across packages.

type Config struct {
	Mods []PathSearch `json:"mods"`
}

type PathSearch struct {
	Path    string       `json:"path"`
	Output  string       `json:"output"`
	Extract string       `json:"extract"`
	Export  string       `json:"export"`
	Include []Include    `json:"include"`
	Exclude [][]string   `json:"exclude"`
	Expects []Expect     `json:"expects"`
	Copy    []PathCopy   `json:"copy"`
	Rename  []PathRename `json:"rename"`
}

type PathRename struct {
	Path []string `json:"path"`
	From []string `json:"from"`
	To   []string `json:"to"`
}

type PathCopy struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Include struct {
	Path string `json:"path"`
	To   string `json:"to"`
}

type Expect struct {
	Path      []string `json:"path"`
	Require   []string `json:"require"`
	Exclusive bool     `json:"exclusive"`
	Base      int      `json:"base"`
}

// Read reads the config file at path and returns a Config.
func Read(path string) (Config, error) {
	data, err := filesystem.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	c := Config{} //nolint:exhaustruct // reason: umarshalling data into struct.
	if err := json.Unmarshal(jsonc.ToJSON(data), &c); err != nil {
		return Config{}, err
	}

	return c, nil
}

// Write writes the config file at path.
func Write(path string, config Config) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	if err := filesystem.WriteFile(path, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func Default() Config {
	return Config{
		Mods: []PathSearch{
			{
				Path:    "pd2mm/pd2/mods",
				Output:  "pd2mm/pd2/output/mods",
				Extract: "pd2mm/pd2/extract/mods",
				Export:  "",
				Include: []Include{},
				Exclude: [][]string{},
				Expects: []Expect{
					{
						Path:      []string{"mod.txt"},
						Require:   []string{},
						Exclusive: false,
						Base:      0,
					},
					{
						Path:      []string{"main.xml"},
						Require:   []string{},
						Exclusive: false,
						Base:      0,
					},
				},
				Copy:   []PathCopy{},
				Rename: []PathRename{},
			},
			{
				Path:    "pd2mm/pd2/mod_overrides",
				Output:  "pd2mm/pd2/output/mod_overrides",
				Extract: "pd2mm/pd2/extract/mod_overrides",
				Export:  "",
				Include: []Include{},
				Exclude: [][]string{},
				Expects: []Expect{
					{Path: []string{"main.xml"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"add.xml"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"effects"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"assets"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"units"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"hooks"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"guis"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"anims"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"soundbanks"}, Require: []string{}, Exclusive: false, Base: 0},
					{Path: []string{"fonts"}, Require: []string{}, Exclusive: false, Base: 0},
				},
				Copy:   []PathCopy{},
				Rename: []PathRename{},
			},
			{
				Path:    "pd2mm/pd2/maps",
				Output:  "pd2mm/pd2/output/maps",
				Extract: "pd2mm/pd2/extract/maps",
				Export:  "",
				Include: []Include{},
				Exclude: [][]string{},
				Expects: []Expect{
					{
						Path:      []string{"main.xml"},
						Require:   []string{},
						Exclusive: false,
						Base:      0,
					},
				},
				Copy:   []PathCopy{},
				Rename: []PathRename{},
			},
		},
	}
}
