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

//nolint:gochecknoglobals // reason: language strings.
package lang

var En = map[string]string{
	// General
	"programName":              "pd2mm",
	"versionUsage":             "The program version",
	"configUsage":              "The config file path",
	"languageNotFound":         "Language not found",
	"extractingNotify":         "... EXTRACTING",
	"copyingNotify":            "... COPYING",
	"startingRunnerNotify":     "... [RUNNER] STARTING",
	"startingCleanerNotify":    "... [CLEANER] STARTING",
	"doneRunnerNotify":         "... [RUNNER] DONE",
	"doneCleanerNotify":        "... [CLEANER] DONE",
	"workingNotify":            "... WORKING",
	"deleteNotify":             "... DELETING",
	"extractNotify":            "... EXTRACTING",
	"doneExtractCleanerNotify": "... EXTRACT CLEANER DONE ...",
	"doneExportCleanerNotify":  "... EXPORT CLEANER DONE ...",
	"doneOutputCleanerNotify":  "... OUTPUT CLEANER DONE ...",
	"errorNotify":              "ERROR:",
	"defaultConfigPath":        "pd2mm/pd2.json",
	"defaultLogPath":           "pd2mm_log.txt",
	"watermarkPart1":           "This work is free of charge",
	"watermarkPart2":           "If you paid money, you were scammed",

	"configLabel":        "Select from available configs",
	"configCustomLabel":  "Set a custom config path",
	"logLabel":           "The log file path",
	"binLabel":           "The 7z file path",
	"startButton":        "Start",
	"cleanExtractButton": "Clean Extract Directories",
	"cleanExportButton":  "Clean Export Directories",
	"cleanOutputButton":  "Clean Output Directories",
}
