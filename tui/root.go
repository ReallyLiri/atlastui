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

type modelState struct {
	selectedSchema string
	selectedTable  string
	tableSection   TableDetailsSection
	quitting       bool
	width          int
	height         int
}

type modelConfig struct {
	title  string
	keymap help.KeyMap
}

type viewModels struct {
	help         help.Model
	tablesList   table.Model
	tableDetails map[TableDetailsSection]table.Model
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
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal size: %w", err)
	}

	m := &model{
		schemasByName:         make(map[string]inspect.Schema),
		tablesBySchemaAndName: make(map[tableKey]inspect.Table),
		state: modelState{
			tableSection: ColumnsTable,
			width:        width,
			height:       height,
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
	tablesListWidth, tablesListHeight := m.tablesListSize()
	m.vms.tablesList = newTablesList(lo.Map(m.schemasByName[m.state.selectedSchema].Tables, func(tbl inspect.Table, _ int) string {
		return tbl.Name
	}), tablesListWidth, tablesListHeight)
	m.vms.tableDetails = nil
}

func (m *model) tablesListSize() (width, height int) {
	return m.state.width * 1 / 3, m.state.height - 5
}

func (m *model) tableDetailsSize() (width, height int) {
	return m.state.width * 2 / 3, m.state.height - 5
}

func (m *model) onTableSelected(table tableKey) {
	t := m.tablesBySchemaAndName[table]
	m.state.selectedTable = t.Name
	w, h := m.tableDetailsSize()
	m.vms.tableDetails = newTableDetails(t, w, h)
}

func (m *model) onResize(width, height int) {
	m.state.width = width
	m.state.height = height

	tablesListWidth, tablesListHeight := m.tablesListSize()
	m.vms.tablesList.SetWidth(tablesListWidth)
	m.vms.tablesList.SetHeight(tablesListHeight)
	tableDetailsWidth, tableDetailsHeight := m.tableDetailsSize()

	for _, t := range m.vms.tableDetails {
		t.SetWidth(tableDetailsWidth)
		t.SetHeight(tableDetailsHeight)
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := msg.(type) {
	case tea.WindowSizeMsg:
		if tmsg.Width > 0 && tmsg.Height > 0 {
			m.onResize(tmsg.Width, tmsg.Height)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(tmsg, keymap.Left), key.Matches(tmsg, keymap.Right):
			if m.vms.tablesList.Focused() {
				m.vms.tablesList.Blur()
				currSection := m.vms.tableDetails[m.state.tableSection]
				currSection.Focus()
				m.vms.tableDetails[m.state.tableSection] = currSection
			} else {
				m.vms.tablesList.Focus()
				currSection := m.vms.tableDetails[m.state.tableSection]
				currSection.Blur()
				m.vms.tableDetails[m.state.tableSection] = currSection
			}
		case key.Matches(tmsg, keymap.Up), key.Matches(tmsg, keymap.Down):
			if m.vms.tablesList.Focused() {
				m.vms.tablesList, cmd = m.vms.tablesList.Update(msg)
				return m, cmd
			}
			if m.vms.tableDetails[m.state.tableSection].Focused() {
				m.vms.tableDetails[m.state.tableSection], cmd = m.vms.tableDetails[m.state.tableSection].Update(msg)
				return m, cmd
			}
		case key.Matches(tmsg, keymap.Help):
			m.vms.help.ShowAll = !m.vms.help.ShowAll
		case key.Matches(tmsg, keymap.Quit):
			m.state.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string {
	if m.state.quitting {
		return ""
	}

	title := titleView(m.config.title, m.state.selectedSchema, m.state.selectedTable)
	contents := []string{
		title,
		strings.Repeat("-", lipgloss.Width(title)),
	}
	if m.state.selectedSchema != "" {
		if m.state.selectedTable != "" {
			contents = append(contents, lipgloss.JoinHorizontal(lipgloss.Left, m.vms.tablesList.View(), " ", m.vms.tableDetails[m.state.tableSection].View()))
		} else {
			contents = append(contents, m.vms.tablesList.View())
		}
	}
	contents = append(contents,
		strings.Repeat("-", lipgloss.Width(title)),
		m.vms.help.View(m.config.keymap),
	)
	return lipgloss.JoinVertical(lipgloss.Top, contents...)
}
