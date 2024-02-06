package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/reallyliri/atlastui/inspect"
	"github.com/samber/lo"
	"strings"
)

type TableDetailsSection string

const (
	ColumnsTable     TableDetailsSection = "Columns"
	IndexesTable     TableDetailsSection = "Indexes"
	ForeignKeysTable TableDetailsSection = "Foreign Keys"
)

func newTableDetails(t inspect.Table, width, height int) map[TableDetailsSection]table.Model {
	tables := make(map[TableDetailsSection]table.Model, 3)
	hasIndexes := len(t.Indexes) > 0
	hasForeignKeys := len(t.ForeignKeys) > 0

	tables[ColumnsTable] = table.New(
		table.WithColumns([]table.Column{
			{Title: "Name", Width: width / 2},
			{Title: "Type", Width: width / 3},
			{Title: "Null", Width: width / 6},
		}),
		table.WithRows(lo.Map(t.Columns, func(col inspect.Column, _ int) table.Row {
			return table.Row{
				col.Name,
				col.Type,
				formatBool(col.Null),
			}
		})),
		table.WithFocused(true),
		table.WithWidth(width),
		table.WithHeight(height),
	)
	if hasIndexes {
		tables[IndexesTable] = table.New(
			table.WithColumns([]table.Column{
				{Title: "Name", Width: 5 * width / 9},
				{Title: "Unique", Width: width / 9},
				{Title: "Parts", Width: 3 * width / 9},
			}),
			table.WithRows(lo.Map(t.Indexes, func(idx inspect.Index, _ int) table.Row {
				return table.Row{
					idx.Name,
					formatBool(idx.Unique),
					strings.Join(lo.Map(idx.Parts, func(part inspect.IndexPart, _ int) string {
						return part.Column
					}), inlineListSeparator),
				}
			})),
			table.WithWidth(width),
			table.WithHeight(height),
		)
	}
	if hasForeignKeys {
		tables[ForeignKeysTable] = table.New(
			table.WithColumns([]table.Column{
				{Title: "Name", Width: width / 2},
				{Title: "Columns", Width: width / 4},
				{Title: "References", Width: width / 4},
			}),
			table.WithRows(lo.Map(t.ForeignKeys, func(fk inspect.ForeignKey, _ int) table.Row {
				return table.Row{
					fk.Name,
					strings.Join(fk.Columns, inlineListSeparator),
					fmt.Sprintf("%s(%s)", fk.References.Table, strings.Join(fk.References.Columns, inlineListSeparator)),
				}
			})),
			table.WithWidth(width),
			table.WithHeight(height),
		)
	}
	return tables
}
