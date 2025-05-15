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
	"io"
	"os"

	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/util"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/lang"
)

// StartConsoleApp is the main entry point for pd2mm.
func StartConsoleApp(logFile io.Writer, version func()) {
	logger.SharedLogger = logger.NewMultiLogger(logFile, os.Stdout)

	errCh := make(chan error, 3) //nolint:mnd // reason: max errors in channel.

	util.DrawWatermark([]string{lang.Lang("programName"), lang.Lang("watermarkPart1"), lang.Lang("watermarkPart2")}, func(s string) {
		logger.SharedLogger.Info(s)
	})

	if data.Flag.Version {
		version()
		return
	}

	if util.IsFlagPassed("config") && data.Flag.Config == "" {
		logger.SharedLogger.Error("flag 'config' cannot be nil or empty")
		return
	}

	if !util.IsFlagPassed("config") {
		data.Flag.Config = ""
	}

	Flags{Flags: data.Flag}.RunWithError(Configs(Flags{Flags: data.Flag}), errCh)

	for err := range errCh {
		if err != nil {
			logger.SharedLogger.Errorf("%s %v", lang.Lang("errorNotify"), err)
		}
	}
}
