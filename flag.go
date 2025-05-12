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
	"flag"

	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/cli"
	"github.com/hkmh223/pd2mm/internal/lang"
)

var (
	flags    *cli.Flags = NewFlags() //nolint:gochecknoglobals // allowed
	defaults            = cli.Flags{ //nolint:gochecknoglobals // allowed
		Version: false,
		Config:  "pd2mm/config.json",
		Lang:    "en",
	}
)

func NewFlags() *cli.Flags {
	return &defaults
}

//nolint:gochecknoinits // allowed
func init() {
	flag.BoolVar(&flags.Version, "version", defaults.Version, lang.Lang("versionUsage"))
	flag.StringVar(&flags.Config, "config", defaults.Config, lang.Lang("configUsage"))

	if flags.Lang != "" {
		err := lang.SetLanguage(flags.Lang)
		if err != nil {
			logger.SharedLogger.Info(lang.Lang("languageNotFound"))
		}
	}

	flag.Parse()
}
