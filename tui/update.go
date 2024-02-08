package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/atlastui/tui/keymap"
)

const maxWidth = 250

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := msg.(type) {
	case tea.WindowSizeMsg:
		if tmsg.Width > 0 && tmsg.Height > 0 {
			m.state.termWidth = tmsg.Width
			if m.state.termWidth > maxWidth {
				m.state.termWidth = maxWidth
			}
			m.state.termHeight = tmsg.Height
		}
	case spinner.TickMsg:
		m.vms.globe, cmd = m.vms.globe.Update(msg)
	case tea.KeyMsg:
		if m.state.easteregg {
			m.state.easteregg = !m.state.easteregg
			break
		}
		switch {
		case key.Matches(tmsg, keymap.Tab):
			m.state.focused = (m.state.focused + 1) % 3
		case key.Matches(tmsg, keymap.Left), key.Matches(tmsg, keymap.Right), key.Matches(tmsg, keymap.Up), key.Matches(tmsg, keymap.Down):
			switch m.state.focused {
			case tablesListFocused:
				m.vms.tablesList, cmd = m.vms.tablesList.Update(msg)
				m.onTableSelected(tableKey{m.state.selectedSchema, m.vms.tablesList.SelectedItem().FilterValue()})
			case detailsTabFocused:
				if key.Matches(tmsg, keymap.Left) || key.Matches(tmsg, keymap.Right) {
					m.state.selectedTab = (m.state.selectedTab + 1) % 3
				}
			case detailsContentsFocused:
				switch m.state.selectedTab {
				case ColumnsTable:
					m.vms.colsChart, cmd = m.vms.colsChart.Update(msg)
				case IndexesTable:
					m.vms.idxChart, cmd = m.vms.idxChart.Update(msg)
				case ForeignKeysTable:
					m.vms.fksChart, cmd = m.vms.fksChart.Update(msg)
				}
			}
		case key.Matches(tmsg, keymap.Help):
			m.vms.help.ShowAll = !m.vms.help.ShowAll
		case key.Matches(tmsg, keymap.Quit):
			m.state.quitting = true
			return m, tea.Quit
		case tmsg.String() == "a":
			m.state.easteregg = !m.state.easteregg
		}
	}
	return m, cmd
}
