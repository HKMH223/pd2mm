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
	"strings"

	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/lang"
	"github.com/hkmh223/pd2mm/internal/pd2mm"
)

var (
	gitHash   string //nolint:gochecknoglobals // reason: build string.
	buildDate string //nolint:gochecknoglobals // reason: build string.
	buildOn   string //nolint:gochecknoglobals // reason: build string.
)

func version() {
	logger.SharedLogger.Info("version", "go", strings.TrimPrefix(buildOn, "go version "), "revision", gitHash, "date", buildDate)
}

func main() {
	if err := lang.SetupLanguage(); err != nil {
		logger.SharedLogger.Fatal(err)
	}

	data.SetupFlags()

	logFile := pd2mm.OpenLogFile(*data.Flag)
	defer func() {
		if err := logFile.Close(); err != nil {
			logger.SharedLogger.Fatal(err)
		}
	}()

	pd2mm.Setup()
	pd2mm.StartConsoleApp(logFile, version)
}
