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
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/pd2mm"
)

var (
	gitHash   string //nolint:gochecknoglobals // allowed
	buildDate string //nolint:gochecknoglobals // allowed
	buildOn   string //nolint:gochecknoglobals // allowed
)

func version() {
	logger.SharedLogger.Info("version", "go", strings.TrimPrefix(buildOn, "go version "), "revision", gitHash, "date", buildDate)
}

func openLog() *os.File { //nolint:mnd // allowed
	file, err := os.OpenFile(filesystem.GetRelativePath("pd2mm_log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func main() {
	file := openLog()
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	logger.SharedLogger = logger.NewMultiLogger(file, os.Stdout)

	util.DrawWatermark([]string{"pd2mm"}, func(s string) {
		logger.SharedLogger.Info(s)
	})

	if flags.Version {
		version()
		return
	}

	if flags.Config == "" {
		logger.SharedLogger.Error("Flag 'config' cannot be nil or empty")
		return
	}

	if !filesystem.Exists(flags.Config) {
		if err := pd2mm.Write(flags.Config); err != nil {
			logger.SharedLogger.Error("Failed to write configuration file", "err", err)
		}
	}

	if c, err := pd2mm.Read(flags.Config); err == nil {
		c.Start()
	} else {
		logger.SharedLogger.Error("Failed to read configuration file", "err", err)
	}
}
