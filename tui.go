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
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hkmh223/pd2mm/common/logger"
	"github.com/hkmh223/pd2mm/internal/pd2mm"
)

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn // allowed
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type { //nolint:exhaustive // allowed
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			run(m.textInput.Value())
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Input configuration file path:\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func initialModel() model {
	input := textinput.New()
	input.Placeholder = "config.json"
	input.Focus()
	input.CharLimit = 256
	input.Width = 20

	return model{
		textInput: input,
		err:       nil,
	}
}

func run(path string) {
	if path == "" {
		logger.SharedLogger.Fatal("Flag 'config' cannot be nil or empty")
	}

	if c, err := pd2mm.Read(path); err == nil {
		c.Start()
	} else {
		logger.SharedLogger.Fatal("Failed to read configuration file", "err", err)
	}
}
