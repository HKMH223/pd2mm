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
	"slices"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/lang"
)

func Start(flags Flags, configs []Config, update func()) {
	flags.Run(configs, update)
}

func ConfigNames(flags Flags) []string {
	var entries []string

	if flags.Config != "" {
		entries = append(entries, flags.Config)
	} else {
		files, err := filesystem.GetTopFiles(lang.Lang("programName"))
		if err != nil {
			logger.SharedLogger.Fatal("Failed to get files", "err", err)
		}

		for _, file := range files {
			if ext := filesystem.GetFileExtension(file); slices.Contains(data.FileTypes, ext) {
				entries = append(entries, filesystem.FromCwd(lang.Lang("programName"), file))
			}
		}
	}

	return entries
}

func Configs(flags Flags) []Config {
	var configs []Config

	entries := ConfigNames(flags)

	for _, entry := range entries {
		if c, err := data.Read(entry); err == nil {
			configs = append(configs, Config{Config: &c})
		}
	}

	return configs
}
