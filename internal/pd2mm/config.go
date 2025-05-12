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

package pd2mm

import (
	"encoding/json"
	"os"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/tidwall/jsonc"
)

type Config struct {
	Mods []PathSearch `json:"mods"`
}

type PathSearch struct {
	Path    string       `json:"path"`
	Output  string       `json:"output"`
	Extract string       `json:"extract"`
	Export  string       `json:"export"`
	Include []Include    `json:"include"`
	Exclude []string     `json:"exclude"`
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
	Path    []string `json:"path"`
	Require []string `json:"require"`
	Base    int      `json:"base"`
}

// Read reads the config file at path and returns a Config.
func Read(path string) (Config, error) {
	data, err := filesystem.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	c := Config{} //nolint:exhaustruct // allowed
	if err := json.Unmarshal(jsonc.ToJSON(data), &c); err != nil {
		return Config{}, err
	}

	return c, nil
}

// Write writes the config file at path.
func Write(path string) error {
	data, err := json.Marshal(Config{}) //nolint:exhaustruct // allowed
	if err != nil {
		return err
	}

	if err := filesystem.WriteFile(path, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}
