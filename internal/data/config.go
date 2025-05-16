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
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/tidwall/jsonc"
)

var FileTypes = []string{".jsonc", ".json"} //nolint:gochecknoglobals // reason: file types are needed across packages.

type Config struct {
	Mods []PathSearch `json:"mods"`
}

type PathSearch struct {
	Mods    string       `json:"mods"`
	Output  PathInfo     `json:"output"`
	Extract PathInfo     `json:"extract"`
	Export  PathInfo     `json:"export"`
	Include []Include    `json:"include"`
	Exclude [][]string   `json:"exclude"`
	Expects []Expect     `json:"expects"`
	Copy    []PathCopy   `json:"copy"`
	Rename  []PathRename `json:"rename"`
}

type PathInfo struct {
	Path         string     `json:"path"`
	ExcludeClean [][]string `json:"excludeClean"`
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

// Format a slice of paths using the current PathSearch settings.
func (search PathSearch) FormatSlice(slice []string) []string {
	result := []string{}

	for _, item := range slice {
		result = append(result, search.FormatString(item))
	}

	return result
}

// Replace keywords with relevant PathSearch settings.
func (search PathSearch) FormatString(str string) string {
	return util.Format(str, map[string]string{
		"{path}":    search.Mods,
		"{output}":  search.Output.Path,
		"{extract}": search.Extract.Path,
		"{export}":  search.Export.Path,
	})
}

//nolint:funlen // reason: setting the default config
func Default() Config {
	return Config{
		Mods: []PathSearch{
			{
				Mods: "pd2mm/pd2/mods",
				Output: PathInfo{
					Path: "pd2mm/pd2/output/mods",
					ExcludeClean: [][]string{
						{"{export}", "saves"},
						{"{export}", "logs"},
						{"{output}", "saves"},
						{"{output}", "logs"},
					},
				},
				Extract: PathInfo{
					Path: "pd2mm/pd2/extract/mods",
					ExcludeClean: [][]string{
						{"{export}", "saves"},
						{"{export}", "logs"},
						{"{output}", "saves"},
						{"{output}", "logs"},
					},
				},
				Export: PathInfo{
					Path: "",
					ExcludeClean: [][]string{
						{"{export}", "saves"},
						{"{export}", "logs"},
						{"{output}", "saves"},
						{"{output}", "logs"},
					},
				},
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
				Mods: "pd2mm/pd2/mod_overrides",
				Output: PathInfo{
					Path:         "pd2mm/pd2/output/mod_overrides",
					ExcludeClean: [][]string{},
				},
				Extract: PathInfo{
					Path:         "pd2mm/pd2/extract/mod_overrides",
					ExcludeClean: [][]string{},
				},
				Export: PathInfo{
					Path:         "",
					ExcludeClean: [][]string{},
				},
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
				Mods: "pd2mm/pd2/mod_overrides",
				Output: PathInfo{
					Path:         "pd2mm/pd2/output/mod_overrides",
					ExcludeClean: [][]string{},
				},
				Extract: PathInfo{
					Path:         "pd2mm/pd2/extract/mod_overrides",
					ExcludeClean: [][]string{},
				},
				Export: PathInfo{
					Path:         "",
					ExcludeClean: [][]string{},
				},
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
