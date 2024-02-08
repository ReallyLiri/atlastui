package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	chart "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/atlastui/inspect"
	"github.com/reallyliri/atlastui/tui/format"
	"github.com/reallyliri/atlastui/tui/keymap"
	"github.com/reallyliri/atlastui/tui/types"
	"github.com/samber/lo"
	"strings"
)

/*
--- NOTE ---
To make the terminology a bit simpler, note "chart" is a table tui component while "table" is a database table.
--- NOTE ---
*/

type tableKey struct {
	schemaName string
	tableName  string
}

type modelState struct {
	selectedSchema string
	selectedTable  string
	selectedTab    types.TableDetailsSection
	focused        types.FocusedComponent
	quitting       bool
	easteregg      bool
	termWidth      int
	termHeight     int
}

type modelConfig struct {
	title  string
	keymap help.KeyMap
}

type viewModels struct {
	help       help.Model
	tablesList list.Model
	colsChart  chart.Model
	idxChart   chart.Model
	fksChart   chart.Model
	globe      spinner.Model
}

type model struct {
	schemasByName         map[string]inspect.Schema
	tablesBySchemaAndName map[tableKey]inspect.Table

	state  modelState
	config modelConfig
	vms    viewModels
}

var _ tea.Model = &model{}

func (m *model) Init() tea.Cmd {
	return m.vms.globe.Tick
}

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
			selectedTab: types.ColumnsTable,
		},
		config: modelConfig{
			keymap: keymap.GetKeyMap(),
			title:  title,
		},
		vms: viewModels{
			help:  help.New(),
			globe: spinner.New(spinner.WithSpinner(spinner.Globe)),
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
	m.vms.tablesList = newTablesList(lo.Map(m.schemasByName[m.state.selectedSchema].Tables, func(chart inspect.Table, _ int) string {
		return chart.Name
	}))
	m.state.selectedTable = ""
}

func (m *model) onTableSelected(key tableKey) {
	m.state.selectedTable = key.tableName
	m.state.selectedTab = types.ColumnsTable
	m.vms.colsChart, m.vms.idxChart, m.vms.fksChart = newCharts(m.tablesBySchemaAndName[key])
}

func newCharts(t inspect.Table) (colsChart chart.Model, idxChart chart.Model, fksChart chart.Model) {
	colsChart = newChart(
		[]chart.Column{
			{Title: "Name", Width: 3},
			{Title: "Type", Width: 2},
			{Title: "Null", Width: 1},
		},
		lo.Map(t.Columns, func(col inspect.Column, _ int) chart.Row {
			return chart.Row{
				format.ColumnName(t, col),
				col.Type,
				format.Bool(col.Null),
			}
		}))

	idxChart = newChart(
		[]chart.Column{
			{Title: "Name", Width: 5},
			{Title: "Unique", Width: 1},
			{Title: "Parts", Width: 3},
		},
		lo.Map(t.Indexes, func(idx inspect.Index, _ int) chart.Row {
			return chart.Row{
				idx.Name,
				format.Bool(idx.Unique),
				strings.Join(lo.Map(idx.Parts, func(part inspect.IndexPart, _ int) string {
					return part.Column
				}), format.InlineListSeparator),
			}
		}))

	fksChart = newChart(
		[]chart.Column{
			{Title: "Name", Width: 2},
			{Title: "Columns", Width: 1},
			{Title: "References", Width: 1},
		},
		lo.Map(t.ForeignKeys, func(fk inspect.ForeignKey, _ int) chart.Row {
			return chart.Row{
				fk.Name,
				strings.Join(fk.Columns, format.InlineListSeparator),
				fmt.Sprintf("%s(%s)", fk.References.Table, strings.Join(fk.References.Columns, format.InlineListSeparator)),
			}
		}))
	return
}

func newChart(cols []chart.Column, rows []chart.Row) chart.Model {
	return chart.New(
		chart.WithColumns(cols),
		chart.WithRows(rows),
		chart.WithFlexColumnWidth(true),
		chart.WithFocused(true),
	)
}

func newTablesList(names []string) list.Model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	lst := list.New(
		lo.Map(names, func(name string, _ int) list.Item {
			return types.TablesListItem(name)
		}),
		delegate,
		0,
		0,
	)
	lst.SetFilteringEnabled(false)
	lst.SetShowHelp(false)
	lst.SetShowTitle(false)
	lst.SetStatusBarItemName("table", "tables")
	lst.SetShowPagination(false)
	return lst
}
