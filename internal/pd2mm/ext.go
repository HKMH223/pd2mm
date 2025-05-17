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
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/data"
	"github.com/hkmh223/pd2mm/internal/lang"
)

// ext.go should place extensions for third party packages here.
// This file should be kept as small as possible.

type (
	Flags      struct{ *data.Flags }
	Config     struct{ *data.Config }
	PathSearch struct{ *data.PathSearch }
	Cleaner    struct{ *data.Cleaner }
)

const (
	Extract = iota + 1
	Export
	Output
)

var SharedCleaner = Cleaner{data.SharedCleaner} //nolint:gochecknoglobals // reason: used by window

// Clean cleans the specified path for each configuration.
func (c Cleaner) Clean(configs []Config, path int, update func() error) {
	errCh := make(chan error, 1)
	c.CleanWithError(configs, path, update, errCh)

	for err := range errCh {
		if err != nil {
			logger.SharedLogger.Errorf("%s %v", lang.Lang("errorNotify"), err)
		}
	}
}

// Clean cleans the specified path for each configuration.
func (c Cleaner) CleanWithError(configs []Config, path int, update func() error, errCh chan<- error) {
	defer close(errCh)

	for _, config := range configs {
		for _, search := range config.Mods {
			switch path {
			case Extract:
				c.Cleaner.Clean(search, search.Extract)
			case Export:
				c.Cleaner.Clean(search, search.Export)
			case Output:
				c.Cleaner.Clean(search, search.Output)
			}
		}
	}

	if err := update(); err != nil {
		errCh <- err
	}
}
