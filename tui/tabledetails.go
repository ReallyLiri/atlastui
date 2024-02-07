package tui

import (
	"fmt"
	chart "github.com/charmbracelet/bubbles/table"
	"github.com/reallyliri/atlastui/inspect"
	"github.com/samber/lo"
	"strings"
)

type tableDetailsSection int

const (
	ColumnsTable tableDetailsSection = iota
	IndexesTable
	ForeignKeysTable
)

func (section tableDetailsSection) title() string {
	switch section {
	case ColumnsTable:
		return "Columns"
	case IndexesTable:
		return "Indexes"
	case ForeignKeysTable:
		return "Foreign Keys"
	default:
		panic("unknown table details section")
	}
}

func newTblDetails(t inspect.Table) (chart.Model, chart.Model, chart.Model) {
	colsTbl := chart.New(
		chart.WithColumns([]chart.Column{
			{Title: "Name", Width: 3},
			{Title: "Type", Width: 2},
			{Title: "Null", Width: 1},
		}),
		chart.WithRows(lo.Map(t.Columns, func(col inspect.Column, _ int) chart.Row {
			return chart.Row{
				formatColumnName(t, col),
				col.Type,
				formatBool(col.Null),
			}
		})),
		chart.WithFlexColumnWidth(true),
		chart.WithFocused(true),
	)

	idxTbl := chart.New(
		chart.WithColumns([]chart.Column{
			{Title: "Name", Width: 5},
			{Title: "Unique", Width: 1},
			{Title: "Parts", Width: 3},
		}),
		chart.WithRows(lo.Map(t.Indexes, func(idx inspect.Index, _ int) chart.Row {
			return chart.Row{
				idx.Name,
				formatBool(idx.Unique),
				strings.Join(lo.Map(idx.Parts, func(part inspect.IndexPart, _ int) string {
					return part.Column
				}), inlineListSeparator),
			}
		})),
		chart.WithFlexColumnWidth(true),
		chart.WithFocused(true),
	)

	fksTbl := chart.New(
		chart.WithColumns([]chart.Column{
			{Title: "Name", Width: 2},
			{Title: "Columns", Width: 1},
			{Title: "References", Width: 1},
		}),
		chart.WithRows(lo.Map(t.ForeignKeys, func(fk inspect.ForeignKey, _ int) chart.Row {
			return chart.Row{
				fk.Name,
				strings.Join(fk.Columns, inlineListSeparator),
				fmt.Sprintf("%s(%s)", fk.References.Table, strings.Join(fk.References.Columns, inlineListSeparator)),
			}
		})),
		chart.WithFlexColumnWidth(true),
		chart.WithFocused(true),
	)

	return colsTbl, idxTbl, fksTbl
}

func formatColumnName(table inspect.Table, col inspect.Column) string {
	sb := strings.Builder{}
	if table.PrimaryKey != nil && lo.ContainsBy(table.PrimaryKey.Parts, func(part inspect.IndexPart) bool {
		return part.Column == col.Name
	}) {
		sb.WriteString("🔑 ")
	} else {
		sb.WriteString("   ")
	}
	sb.WriteString(col.Name)
	if lo.ContainsBy(table.ForeignKeys, func(fk inspect.ForeignKey) bool {
		return lo.Contains(fk.Columns, col.Name)
	}) {
		sb.WriteString(" 🔗")
	}
	return sb.String()
}
