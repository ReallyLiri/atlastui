package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/reallyliri/atlastui/inspect"
	"github.com/reallyliri/atlastui/tui/keymap"
	"github.com/samber/lo"
	"golang.org/x/term"
	"os"
	"strings"
)

type tableKey struct {
	schemaName string
	tableName  string
}

type rootModel struct {
	// data
	schemasByName         map[string]inspect.Schema
	tablesBySchemaAndName map[tableKey]inspect.Table

	// state
	selectedSchema string
	selectedTable  string
	tableSection   TableDetailsSection
	quitting       bool

	// config
	title  string
	keymap help.KeyMap
	width  int
	height int

	// view models
	help         help.Model
	tablesList   table.Model
	tableDetails map[TableDetailsSection]table.Model
}

var _ tea.Model = &rootModel{}

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

func newRootModel(title string, data inspect.Data) (*rootModel, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal size: %w", err)
	}

	m := &rootModel{
		schemasByName:         make(map[string]inspect.Schema),
		tablesBySchemaAndName: make(map[tableKey]inspect.Table),
		keymap:                keymap.GetKeyMap(),
		title:                 title,
		help:                  help.New(),
		tableSection:          ColumnsTable,
		width:                 width,
		height:                height,
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

func (m *rootModel) onSchemaSelected(schema string) {
	m.selectedSchema = schema
	tablesListWidth, tablesListHeight := m.tablesListSize()
	m.tablesList = newTablesList(lo.Map(m.schemasByName[m.selectedSchema].Tables, func(tbl inspect.Table, _ int) string {
		return tbl.Name
	}), tablesListWidth, tablesListHeight)
	m.tableDetails = nil
}

func (m *rootModel) tablesListSize() (width, height int) {
	return m.width * 1 / 3, m.height - 5
}

func (m *rootModel) tableDetailsSize() (width, height int) {
	return m.width * 2 / 3, m.height - 5
}

func (m *rootModel) onTableSelected(table tableKey) {
	t := m.tablesBySchemaAndName[table]
	m.selectedTable = t.Name
	w, h := m.tableDetailsSize()
	m.tableDetails = newTableDetails(t, w, h)
}

func (m *rootModel) onResize(width, height int) {
	m.width = width
	m.height = height

	tablesListWidth, tablesListHeight := m.tablesListSize()
	m.tablesList.SetWidth(tablesListWidth)
	m.tablesList.SetHeight(tablesListHeight)
	tableDetailsWidth, tableDetailsHeight := m.tableDetailsSize()

	for _, t := range m.tableDetails {
		t.SetWidth(tableDetailsWidth)
		t.SetHeight(tableDetailsHeight)
	}
}

func (m *rootModel) Init() tea.Cmd {
	return nil
}

func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := msg.(type) {
	case tea.WindowSizeMsg:
		if tmsg.Width > 0 && tmsg.Height > 0 {
			m.onResize(tmsg.Width, tmsg.Height)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(tmsg, keymap.Left), key.Matches(tmsg, keymap.Right):
			if m.tablesList.Focused() {
				m.tablesList.Blur()
				currSection := m.tableDetails[m.tableSection]
				currSection.Focus()
				m.tableDetails[m.tableSection] = currSection
			} else {
				m.tablesList.Focus()
				currSection := m.tableDetails[m.tableSection]
				currSection.Blur()
				m.tableDetails[m.tableSection] = currSection
			}
		case key.Matches(tmsg, keymap.Up), key.Matches(tmsg, keymap.Down):
			if m.tablesList.Focused() {
				m.tablesList, cmd = m.tablesList.Update(msg)
				return m, cmd
			}
			if m.tableDetails[m.tableSection].Focused() {
				m.tableDetails[m.tableSection], cmd = m.tableDetails[m.tableSection].Update(msg)
				return m, cmd
			}
		case key.Matches(tmsg, keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(tmsg, keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *rootModel) View() string {
	if m.quitting {
		return ""
	}

	title := titleView(m.title, m.selectedSchema, m.selectedTable)
	contents := []string{
		title,
		strings.Repeat("-", lipgloss.Width(title)),
	}
	if m.selectedSchema != "" {
		if m.selectedTable != "" {
			contents = append(contents, lipgloss.JoinHorizontal(lipgloss.Left, m.tablesList.View(), " ", m.tableDetails[m.tableSection].View()))
		} else {
			contents = append(contents, m.tablesList.View())
		}
	}
	contents = append(contents,
		strings.Repeat("-", lipgloss.Width(title)),
		m.help.View(m.keymap),
	)
	return lipgloss.JoinVertical(lipgloss.Top, contents...)
}
