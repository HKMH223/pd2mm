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

package lang

import (
	"errors"
	"sync"

	"golang.org/x/text/language"
)

var _langmap = map[string]map[string]string{ //nolint:gochecknoglobals // reason: map used to set language and get keys.
	"zh": Zh,
	"en": En,
}

var (
	_lang map[string]string //nolint:gochecknoglobals // reason: current language.
	_lock = sync.RWMutex{}  //nolint:gochecknoglobals // reason: lock used across functions.
)

// SetLanguage sets the language of the program.
func SetLanguage(languge string) error {
	_lock.Lock()
	defer _lock.Unlock()

	tag := language.Make(languge)
	lstr, _ := tag.Base()

	if _langmap[lstr.String()] == nil {
		return ErrNotFind
	}

	_lang = _langmap[lstr.String()]

	return nil
}

// Lang returns the value of a key in the current language.
func Lang(key string) string {
	_lock.RLock()
	defer _lock.RUnlock()

	word, ok := _lang[key]
	if !ok {
		return _langmap["en"][key]
	}

	return word
}

var ErrNotFind = errors.New("ErrNotFind")
