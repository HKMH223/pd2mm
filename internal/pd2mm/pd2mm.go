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
	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/hkmh223/pd2mm/common/logger"
)

// Start starts the extraction and processing of mods.
func (c Config) Start() {
	for _, search := range c.Mods {
		if err := filesystem.DeleteDirectory(filesystem.FromCwd(search.Output)); err != nil {
			logger.SharedLogger.Warn("Failed to delete directory", "path", search.Output, "err", err)
		}
	}

	for _, search := range c.Mods {
		if err := Extract(search); err != nil {
			logger.SharedLogger.Error("Failed to extract mods", "err", err)
			continue
		}

		if err := c.Process(search); err != nil {
			logger.SharedLogger.Error("Failed to process mods", "err", err)
			continue
		}
	}
}
