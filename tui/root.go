package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/atlaspect/inspect"
)

type tableKey struct {
	schemaName string
	tableName  string
}

type rootModel struct {
	schemasByName         map[string]inspect.Schema
	tablesBySchemaAndName map[tableKey]inspect.Table
	selectedSchema        string
	selectedTable         string
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

func msgToKey(msg tea.Msg) string {
	kmsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return ""
	}
	return kmsg.String()
}

func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgToKey(msg) {
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *rootModel) View() string {
	//TODO implement me
	panic("implement me")
}
