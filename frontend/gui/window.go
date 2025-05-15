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
	"bytes"
	"io"

	giu "github.com/AllenDang/giu"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/lang"
	"github.com/hkmh223/pd2mm/internal/pd2mm"
)

var (
	sashPos1 float32      = 500 //nolint:gochecknoglobals // allowed
	sashPos2 float32      = 300 //nolint:gochecknoglobals // allowed
	buf      bytes.Buffer       //nolint:gochecknoglobals // allowed
)

// StartApp is the main entry point for pd2mm.
func StartApp(version string, logFile io.Writer) {
	logger.SharedLogger = logger.NewMultiLogger(logFile, &buf)
	logger.SharedLogger.Info("Initialized!")

	wnd := giu.NewMasterWindow("pd2mm - "+version, 840, 500, 0) //nolint:mnd // allowed
	wnd.Run(window)
}

func start() {
	pd2mm.Start(pd2mm.Flags{Flags: data.Flag}, func() { giu.Update() })
}

func window() {
	giu.SingleWindow().Layout(
		giu.Condition(pd2mm.IsRunning, giu.Label("Working..."), nil),
		giu.SplitLayout(giu.DirectionHorizontal, &sashPos2,
			giu.Layout{
				giu.SplitLayout(giu.DirectionVertical, &sashPos1,
					giu.Layout{
						giu.Style().SetDisabled(pd2mm.IsRunning).To(giu.InputText(&data.Flag.Config).Label(lang.Lang("configLabel"))),
						giu.Style().SetDisabled(pd2mm.IsRunning).To(giu.InputText(&data.Flag.Bin).Label(lang.Lang("binLabel"))),
						giu.Style().SetDisabled(pd2mm.IsRunning).To(giu.InputText(&data.Flag.Log).Label(lang.Lang("logLabel"))),
					},
					giu.Layout{
						// giu.Style().SetDisabled(pd2mm.IsRunning).To(giu.Checkbox("Create zip after download", &flags.Zip)),
						giu.Separator(),
						giu.Row(
							giu.Button(lang.Lang("startButton")).OnClick(start).Disabled(pd2mm.IsRunning),
						),
					},
				),
			}, giu.Layout{
				giu.Child().Layout(
					giu.Label(buf.String()),
				),
			}),
	)
}
