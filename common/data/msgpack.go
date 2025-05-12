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

package data

import (
	"encoding/json"

	"github.com/hkmh223/pd2mm/common/filesystem"
	"github.com/tidwall/gjson"
	"github.com/vmihailenco/msgpack"
)

// Encode a file in the msgpack format.
func MsgpackEncode(value any) ([]byte, error) {
	b, err := msgpack.Marshal(value)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Decode a file in the msgpack format.
func MsgpackDecode(data []byte) (any, error) {
	var bytes any

	err := msgpack.Unmarshal(data, &bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Encode a file in the msgpack format.
func MsgpackEncodeFile(name string) ([]byte, error) {
	b, err := filesystem.ReadFile(name)
	if err != nil {
		return nil, err
	}

	enc, err := MsgpackEncode(gjson.ParseBytes(b).Value())
	if err != nil {
		return nil, err
	}

	return enc, nil
}

// Decode a file in the msgpack format.
func MsgpackDecodeFile(name string) (string, error) {
	b, err := filesystem.ReadFile(name)
	if err != nil {
		return "", err
	}

	dec, err := MsgpackDecode(b)
	if err != nil {
		return "", err
	}

	json, err := json.MarshalIndent(dec, "", " ")
	if err != nil {
		return "", err
	}

	return string(json), nil
}
