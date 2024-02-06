package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/atlastui/inspect"
	"github.com/reallyliri/atlastui/tui/keymap"
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
	quitting       bool

	// config
	keymap help.KeyMap

	// view models
	help help.Model
}

var _ tea.Model = &rootModel{}

func Run(ctx context.Context, data inspect.Data) error {
	m := newRootModel(data)
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

func newRootModel(data inspect.Data) *rootModel {
	m := &rootModel{
		schemasByName:         make(map[string]inspect.Schema),
		tablesBySchemaAndName: make(map[tableKey]inspect.Table),
		keymap:                keymap.GetKeyMap(),
		help:                  help.New(),
	}

	for _, schema := range data.Schemas {
		m.schemasByName[schema.Name] = schema
		for _, table := range schema.Tables {
			m.tablesBySchemaAndName[tableKey{schema.Name, table.Name}] = table
		}
	}

	if len(data.Schemas) == 1 {
		schema := data.Schemas[0]
		m.selectedSchema = schema.Name
		if len(schema.Tables) == 1 {
			m.selectedTable = schema.Tables[0].Name
		}
	}
	return m
}

func (m *rootModel) Init() tea.Cmd {
	return nil
}

func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO - handle window resize
	case tea.KeyMsg:
		switch {
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

	sb := strings.Builder{}
	sb.WriteString("WIP")
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d schemas\n", len(m.schemasByName)))
	sb.WriteString(fmt.Sprintf("%d tables\n", len(m.tablesBySchemaAndName)))
	sb.WriteString(m.help.View(m.keymap))
	return sb.String()
}
