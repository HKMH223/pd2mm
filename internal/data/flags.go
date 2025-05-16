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
	"flag"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/lang"
)

type Flags struct {
	Version      bool
	Config       string
	Log          string
	Lang         string
	Bin          string
	CleanExtract bool
	CleanExport  bool
	CleanOutput  bool
}

var (
	Flag     = NewFlags() //nolint:gochecknoglobals // reason: flags are needed across packages.
	defaults = Flags{     //nolint:gochecknoglobals // reason: default flags are assigned in multiple functions.
		Version:      false,
		Config:       lang.Lang("defaultConfigPath"),
		Log:          lang.Lang("defaultLogPath"),
		Lang:         "en",
		Bin:          filesystem.Normalize(filesystem.Combine(lang.Lang("programName"), "bin")),
		CleanExtract: false,
		CleanExport:  false,
		CleanOutput:  false,
	}
)

// NewFlags creates a new Flags instance.
func NewFlags() *Flags {
	return &defaults
}

//nolint:gochecknoinits // reason: setup program flags.
func init() {
	flag.BoolVar(&Flag.Version, "version", defaults.Version, lang.Lang("versionUsage"))
	flag.StringVar(&Flag.Config, "config", defaults.Config, lang.Lang("configUsage"))

	if Flag.Lang != "" {
		err := lang.SetLanguage(Flag.Lang)
		if err != nil {
			logger.SharedLogger.Info(lang.Lang("languageNotFound"))
		}
	}

	flag.Parse()
}
