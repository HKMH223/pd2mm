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

// Setup creates the default config if it does not exist.
func Setup() {
	path, err := filesystem.FromCwd(lang.Lang("defaultConfigPath"))
	if err != nil {
		logger.SharedLogger.Error(err)
	}

	if !filesystem.Exists(path) {
		logger.SharedLogger.Warnf("%s does not exist, creating.", lang.Lang("defaultConfigPath"))

		if err := data.Write(lang.Lang("defaultConfigPath"), data.Default()); err != nil {
			logger.SharedLogger.Error("failed to write default config", "err", err)
			return
		}
	}
}

// Start starts the program.
func Start(flags Flags, configs []Config, update func() error) {
	SharedRunner.RegisterUpdate(update)
	SharedRunner.Run(flags, configs)
}

// ConfigNames returns the names of the configs.
func ConfigNames(flags Flags) ([]string, error) {
	var entries []string

	if flags.Config != "" {
		entries = append(entries, flags.Config)
	} else {
		files, err := filesystem.GetTopFiles(lang.Lang("programName"))
		if err != nil {
			logger.SharedLogger.Error("failed to get files", "err", err)
			return entries, nil
		}

		for _, file := range files {
			path, err := filesystem.FromCwd(lang.Lang("programName"), file)
			if err != nil {
				return nil, err
			}

			if ext := filesystem.GetFileExtension(file); slices.Contains(data.FileTypes, ext) {
				entries = append(entries, path)
			}
		}
	}

	return entries, nil
}

// Configs returns the configs.
func Configs(flags Flags) ([]Config, error) {
	var configs []Config

	entries, err := ConfigNames(flags)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if c, err := data.Read(entry); err == nil {
			configs = append(configs, Config{Config: &c})
		}
	}

	return configs, nil
}
