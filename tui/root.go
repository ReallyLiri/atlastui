package tui

import (
	"context"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tbl "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/reallyliri/atlastui/inspect"
	"github.com/reallyliri/atlastui/tui/keymap"
	"github.com/samber/lo"
	"strings"
)

/*
To make the terminology a bit simpler, note "tbl" is a table tui component, "table" is a database table.
*/

const maxWidth = 250

type tableKey struct {
	schemaName string
	tableName  string
}

type focusedComponent int

const (
	tablesListFocused focusedComponent = iota
	detailsTabFocused
	detailsContentsFocused
)

type modelState struct {
	selectedSchema string
	selectedTable  string
	selectedTab    tableDetailsSection
	quitting       bool
	termWidth      int
	termHeight     int
	focused        focusedComponent
}

type modelConfig struct {
	title  string
	keymap help.KeyMap
}

type viewModels struct {
	help       help.Model
	tablesList list.Model
	colsTbl    tbl.Model
	idxTbl     tbl.Model
	fksTbl     tbl.Model
}

type model struct {
	schemasByName         map[string]inspect.Schema
	tablesBySchemaAndName map[tableKey]inspect.Table

	state  modelState
	config modelConfig
	vms    viewModels
}

var _ tea.Model = &model{}

func Run(ctx context.Context, title string, data inspect.Data) error {
	m, err := newRootModel(title, data)
	if err != nil {
		return err
	}
	prog := tea.NewProgram(m)
	go func() {
		<-ctx.Done()
		prog.Quit()
	}()
	if _, err := prog.Run(); err != nil {
		return err
	}
	return nil
}

func newRootModel(title string, data inspect.Data) (*model, error) {
	m := &model{
		schemasByName:         make(map[string]inspect.Schema),
		tablesBySchemaAndName: make(map[tableKey]inspect.Table),
		state: modelState{
			selectedTab: ColumnsTable,
		},
		config: modelConfig{
			keymap: keymap.GetKeyMap(),
			title:  title,
		},
		vms: viewModels{
			help: help.New(),
		},
	}

	for _, schema := range data.Schemas {
		m.schemasByName[schema.Name] = schema
		for _, table := range schema.Tables {
			m.tablesBySchemaAndName[tableKey{schema.Name, table.Name}] = table
		}
	}

	// TODO - support multiple schemas
	schema := data.Schemas[0]
	m.onSchemaSelected(schema.Name)
	m.onTableSelected(tableKey{schema.Name, schema.Tables[0].Name})

	return m, nil
}

func (m *model) onSchemaSelected(schema string) {
	m.state.selectedSchema = schema
	m.vms.tablesList = newTablesList(lo.Map(m.schemasByName[m.state.selectedSchema].Tables, func(tbl inspect.Table, _ int) string {
		return tbl.Name
	}))
	m.state.selectedTable = ""
}

func (m *model) onTableSelected(key tableKey) {
	m.state.selectedTable = key.tableName
	m.state.selectedTab = ColumnsTable
	m.vms.colsTbl, m.vms.idxTbl, m.vms.fksTbl = newTblDetails(m.tablesBySchemaAndName[key], 150, 0)
}

func (m *model) Init() tea.Cmd {
	return nil
}

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
	case tea.KeyMsg:
		switch {
		case key.Matches(tmsg, keymap.Left), key.Matches(tmsg, keymap.Right):
			m.state.focused = (m.state.focused + 1) % 3
		case key.Matches(tmsg, keymap.Up), key.Matches(tmsg, keymap.Down):
			switch m.state.focused {
			case tablesListFocused:
				m.vms.tablesList, cmd = m.vms.tablesList.Update(msg)
				m.onTableSelected(tableKey{m.state.selectedSchema, m.vms.tablesList.SelectedItem().FilterValue()})
			case detailsTabFocused:
				m.state.selectedTab = (m.state.selectedTab + 1) % 3
			case detailsContentsFocused:
				switch m.state.selectedTab {
				case ColumnsTable:
					m.vms.colsTbl, cmd = m.vms.colsTbl.Update(msg)
				case IndexesTable:
					m.vms.idxTbl, cmd = m.vms.idxTbl.Update(msg)
				case ForeignKeysTable:
					m.vms.fksTbl, cmd = m.vms.fksTbl.Update(msg)
				}
			}
		case key.Matches(tmsg, keymap.Help):
			m.vms.help.ShowAll = !m.vms.help.ShowAll
		case key.Matches(tmsg, keymap.Quit):
			m.state.quitting = true
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m *model) View() string {
	if m.state.quitting {
		return ""
	}
	borderWidth, borderHeight := borderFocusedStyle.GetFrameSize()

	title := titleView(m.config.title, m.state.selectedSchema, m.state.selectedTable)
	footer := m.vms.help.View(m.config.keymap)
	centerHeight := m.state.termHeight - lipgloss.Height(title) - lipgloss.Height(footer) - 5

	var tablesList string
	var tabsView string
	var details string

	if m.state.selectedSchema != "" {
		m.vms.tablesList.SetSize(m.state.termWidth/3-borderWidth, centerHeight-borderHeight+3)
		tablesList = withBorder(m.vms.tablesList.View(), m.state.focused == tablesListFocused)

		if m.state.selectedTable != "" {
			detailsWidth := (m.state.termWidth*2)/3 - borderWidth
			tabsView = m.tabsView(detailsWidth, m.state.focused == detailsTabFocused)

			var currTbl tbl.Model
			switch m.state.selectedTab {
			case ColumnsTable:
				currTbl = m.vms.colsTbl
			case IndexesTable:
				currTbl = m.vms.idxTbl
			case ForeignKeysTable:
				currTbl = m.vms.fksTbl
			}
			currTbl.SetWidth(detailsWidth)
			currTbl.SetHeight(centerHeight - lipgloss.Height(tabsView) - borderHeight + 2)
			details = withBorder(currTbl.View(), m.state.focused == detailsContentsFocused)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			tablesList,
			lipgloss.JoinVertical(
				lipgloss.Top,
				tabsView, details,
			),
		),
		footer,
	)
}

func (m *model) tabsView(width int, focused bool) string {
	tabs := []string{
		tabView("Columns", m.state.selectedTab == ColumnsTable),
		tabView("Indexes", m.state.selectedTab == IndexesTable),
		tabView("Foreign Keys", m.state.selectedTab == ForeignKeysTable),
	}

	row := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(strings.Join(tabs, subTitleStyle.Render(" · ")))
	return withBorder(row, focused)
}

func tabView(title string, selected bool) string {
	return lo.Ternary(selected, titleStyle.Render(title), subTitleStyle.Render(title))
}

func withBorder(ui string, focused bool) string {
	if focused {
		return borderFocusedStyle.Render(ui)
	} else {
		return borderBluredStyle.Render(ui)
	}
}
