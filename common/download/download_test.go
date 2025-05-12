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

package download_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hkmh223/pd2mm/common/download"
)

const testURL = "https://raw.githubusercontent.com/hkmh223/pd2mm/main/LICENSE"

func TestGenericDownload(t *testing.T) {
	t.Parallel()

	msg := download.Messenger{
		StartDownload: func(fname string) {
			fmt.Printf("Test download: %s\n", fname)
		},
	}

	if bytes, err := download.Download(testURL); err != nil || len(bytes) == 0 {
		t.Fatal("download fail")
	}

	if bytes, err := download.WithContext(context.TODO(), msg, testURL); err != nil || len(bytes) == 0 {
		t.Fatal("download fail")
	}
}

func TestFileDownload(t *testing.T) {
	t.Parallel()

	if err := download.File(testURL, "LICENSE", "./.test/"); err != nil {
		t.Fatal(err)
	}

	if bytes, err := download.FileWithBytes(testURL, "LICENSE", "./.test/"); err != nil || len(bytes) == 0 {
		t.Fatal("download fail")
	}
}

//nolint:lll // test only
func TestFileValidated(t *testing.T) {
	t.Parallel()

	if err := download.FileValidated(testURL, "aaabbbccc", "LICENSE", "./.test/"); err == nil {
		t.Fatal("download fail")
	}

	if bytes, err := download.FileWithBytesValidated(testURL, "aaabbbccc", "LICENSE", "./.test/"); err == nil || len(bytes) != 0 {
		t.Fatal("download fail")
	}

	if err := download.FileValidated(testURL, "8486a10c4393cee1c25392769ddd3b2d6c242d6ec7928e1414efff7dfb2f07ef", "LICENSE", "./.test/"); err != nil {
		t.Fatal(err)
	}

	if bytes, err := download.FileWithBytesValidated(testURL, "8486a10c4393cee1c25392769ddd3b2d6c242d6ec7928e1414efff7dfb2f07ef", "LICENSE", "./.test/"); err != nil || len(bytes) == 0 {
		t.Fatal("download fail")
	}
}

//nolint:lll // test only
func TestFileDownloadWithHash(t *testing.T) {
	t.Parallel()

	msg := download.Messenger{
		StartDownload: func(fname string) {
			fmt.Printf("Test download: %s\n", fname)
		},
	}

	if err := download.FileWithContext(context.TODO(), msg, testURL, "8486a10c4393cee1c25392769ddd3b2d6c242d6ec7928e1414efff7dfb2f07ef", "LICENSE", "./.test/", download.DefaultHashValidator); err != nil {
		t.Fatal(err)
	}

	if err := download.FileWithContext(context.TODO(), msg, testURL, "", "LICENSE", "./.test/", download.DefaultHashValidator); err == nil {
		t.Fatal("empty hash has validated successfully")
	}

	if bytes, err := download.FileWithContextAndBytes(context.TODO(), msg, testURL, "", "LICENSE", "./.test/", nil); err != nil || len(bytes) == 0 {
		t.Fatal("download fail")
	}
}
