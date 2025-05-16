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

	giu "github.com/AllenDang/giu"
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/common/safe"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/lang"
	"github.com/hkmh223/pd2mm/internal/pd2mm"
)

//nolint:gochecknoglobals // reason: used by multiple functions.
var (
	_width                  = 840
	_height                 = 500
	_sashPos1       float32 = 500
	_sashPos2       float32 = 300
	_buf                    = filesystem.NewLineRingBuffer(100) //nolint:mnd // reason: line count
	_configs        []string
	_selectedConfig int32
	_disabled       bool
)

// StartApp is the main entry point for pd2mm.
func StartApp(version string, logFile io.Writer) error {
	logger.RegisterLogger(logFile, _buf)
	logger.SharedLogger.Info("Initialized!")

	// ConfigNames either takes a the flag config, otherwise get all configs in the directory.
	// We want to get all configs.
	data.Flag.Config = ""

	var err error

	_configs, err = pd2mm.ConfigNames(pd2mm.Flags{Flags: data.Flag})
	if err != nil {
		return err
	}

	if len(_configs) == 0 {
		_configs = append(_configs, lang.Lang("defaultConfigPath"))
	}

	wnd := giu.NewMasterWindow("pd2mm - "+version, _width, _height, 0)
	wnd.Run(window)

	return nil
}

// startButton is the button that starts pd2mm.
func startButton() {
	configs, err := readConfigs()
	if err != nil {
		logger.SharedLogger.Error("failed to read configuration file", "err", err)

		return
	}

	//nolint:unparam // reason: update does not return error
	go pd2mm.Start(pd2mm.Flags{Flags: data.Flag}, configs, func() error {
		giu.Update()
		return nil
	})
}

// cleanExtractDirectoryButton is the button that cleans the extract path.
func cleanExtractDirectoryButton() {
	configs, err := readConfigs()
	if err != nil {
		logger.SharedLogger.Error("failed to read configuration file", "err", err)

		return
	}

	//nolint:unparam // reason: update does not return error
	go pd2mm.SharedCleaner.Clean(configs, pd2mm.Extract, func() error {
		giu.Update()
		return nil
	})
}

// cleanExportDirectoryButton is the button that cleans the export path.
func cleanExportDirectoryButton() {
	configs, err := readConfigs()
	if err != nil {
		logger.SharedLogger.Error("failed to read configuration file", "err", err)

		return
	}

	//nolint:unparam // reason: update does not return error
	go pd2mm.SharedCleaner.Clean(configs, pd2mm.Export, func() error {
		giu.Update()
		return nil
	})
}

// CleanOutputButton is the button that cleans the output path.
func cleanOutputDirectoryButton() {
	configs, err := readConfigs()
	if err != nil {
		logger.SharedLogger.Error("failed to read configuration file", "err", err)

		return
	}

	//nolint:unparam // reason: update does not return error
	go pd2mm.SharedCleaner.Clean(configs, pd2mm.Output, func() error {
		giu.Update()
		return nil
	})
}

// Read all configs and return a slice of pd2mm.Configs.
// Generally reading configs every time you need them isn't great, you could load them all once on startup.
// However, it makes debugging capabilities much easier.
func readConfigs() ([]pd2mm.Config, error) {
	config, err := data.Read(safe.Slice(_configs, int(_selectedConfig)))
	if err != nil {
		return nil, err
	}

	configs := []pd2mm.Config{{Config: &config}}
	if data.Flag.Config != "" {
		configs, err = pd2mm.Configs(pd2mm.Flags{Flags: data.Flag})
		if err != nil {
			return nil, err
		}
	}

	return configs, nil
}

//nolint:lll // reason: function chaining is used by giu.
func window() {
	if pd2mm.SharedRunner.IsActive() || pd2mm.SharedCleaner.IsActive() {
		_disabled = true
	} else if !pd2mm.SharedRunner.IsActive() && !pd2mm.SharedCleaner.IsActive() {
		_disabled = false
	}

	giu.SingleWindow().Layout(
		giu.Condition(_disabled, giu.Label(lang.Lang("workingNotify")), nil),
		giu.SplitLayout(giu.DirectionHorizontal, &_sashPos2,
			giu.Layout{
				giu.SplitLayout(giu.DirectionVertical, &_sashPos1,
					giu.Layout{
						giu.Style().SetDisabled(_disabled).To(giu.InputText(&data.Flag.Bin).Label(lang.Lang("binLabel"))),
						giu.Style().SetDisabled(_disabled).To(giu.InputText(&data.Flag.Log).Label(lang.Lang("logLabel"))),
						giu.Style().SetDisabled(_disabled).To(giu.Combo(lang.Lang("configLabel"), safe.Slice(_configs, int(_selectedConfig)), _configs, &_selectedConfig)),
						giu.Style().SetDisabled(_disabled).To(giu.InputText(&data.Flag.Config).Hint(lang.Lang("defaultConfigPath")).Label(lang.Lang("configCustomLabel"))),
					},
					giu.Layout{
						giu.Separator(),
						giu.Column(
							giu.Button(lang.Lang("startButton")).OnClick(startButton).Disabled(_disabled).Size(-1, 0),
							giu.Button(lang.Lang("cleanExtractButton")).OnClick(cleanExtractDirectoryButton).Disabled(_disabled).Size(-1, 0),
							giu.Button(lang.Lang("cleanExportButton")).OnClick(cleanExportDirectoryButton).Disabled(_disabled).Size(-1, 0),
							giu.Button(lang.Lang("cleanOutputButton")).OnClick(cleanOutputDirectoryButton).Disabled(_disabled).Size(-1, 0),
						),
					},
				),
			}, giu.Layout{
				giu.Child().Layout(
					giu.Label(_buf.String()),
				),
			}),
	)
}
