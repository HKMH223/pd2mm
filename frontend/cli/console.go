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
	"io"
	"os"

	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/lang"
	"github.com/hkmh223/pd2mm/internal/pd2mm"
)

// StartApp is the main entry point for pd2mm.
func StartApp(logFile io.Writer) {
	logger.SharedLogger = logger.NewMultiLogger(logFile, os.Stdout)

	util.DrawWatermark([]string{lang.Lang("programName"), lang.Lang("watermarkPart1"), lang.Lang("watermarkPart2")}, func(s string) {
		logger.SharedLogger.Info(s)
	})

	if data.Flag.Version {
		version()
		return
	}

	if util.IsFlagPassed("config") && data.Flag.Config == "" {
		logger.SharedLogger.Fatal("Flag 'config' cannot be nil or empty")
	}

	pd2mm.Start(pd2mm.Flags{Flags: data.Flag}, func() {})
}
