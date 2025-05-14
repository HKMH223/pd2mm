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
	logFile := openLog()
	defer func() {
		if err := logFile.Close(); err != nil {
			panic(err)
		}
	}()

	StartApp(logFile)
}
