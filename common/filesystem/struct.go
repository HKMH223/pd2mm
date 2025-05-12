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

package filesystem

type PathCheck struct {
	Type   string
	Target string
	Action string
}

const (
	PathCheckTypeEndsWith  string = "EndsWith"
	PathCheckTypeContains  string = "Contains"
	PathCheckTypeDriveRoot string = "DriveRoot"

	PathCheckActionWarn string = "Warn"
	PathCheckActionDeny string = "Deny"
)

// Default list of problematic paths.
func DefaultProblemPaths() []PathCheck {
	return []PathCheck{
		{Type: PathCheckTypeEndsWith, Target: "SteamApps", Action: PathCheckActionWarn},
		{Type: PathCheckTypeEndsWith, Target: "Documents", Action: PathCheckActionWarn},
		{Type: PathCheckTypeEndsWith, Target: "Desktop", Action: PathCheckActionDeny},
		{Type: PathCheckTypeContains, Target: "Desktop", Action: PathCheckActionWarn},
		{Type: PathCheckTypeContains, Target: "scoped_dir", Action: PathCheckActionDeny},
		{Type: PathCheckTypeContains, Target: "Downloads", Action: PathCheckActionDeny},
		{Type: PathCheckTypeContains, Target: "OneDrive", Action: PathCheckActionDeny},
		{Type: PathCheckTypeContains, Target: "NextCloud", Action: PathCheckActionDeny},
		{Type: PathCheckTypeContains, Target: "DropBox", Action: PathCheckActionDeny},
		{Type: PathCheckTypeContains, Target: "Google", Action: PathCheckActionDeny},
		{Type: PathCheckTypeContains, Target: "Program Files", Action: PathCheckActionDeny},
		{Type: PathCheckTypeContains, Target: "Program Files (x86)", Action: PathCheckActionDeny},
		// {Type: PathCheckTypeContains, Target: "Windows", Action: PathCheckActionDeny},
		{Type: PathCheckTypeDriveRoot, Target: "", Action: PathCheckActionDeny},

		// Reserved words
		{Type: PathCheckTypeEndsWith, Target: "CON", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "PRN", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "AUX", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "CLOCK$", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "NUL", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM0", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM1", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM2", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM3", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM4", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM5", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM6", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM7", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM8", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "COM9", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT0", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT1", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT2", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT3", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT4", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT5", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT6", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT7", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT8", Action: PathCheckActionDeny},
		{Type: PathCheckTypeEndsWith, Target: "LPT9", Action: PathCheckActionDeny},
	}
}
