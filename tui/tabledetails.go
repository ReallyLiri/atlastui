package tui

import (
	"fmt"
	tbl "github.com/charmbracelet/bubbles/table"
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

func newTblDetails(t inspect.Table) (tbl.Model, tbl.Model, tbl.Model) {
	colsTbl := tbl.New(
		tbl.WithColumns([]tbl.Column{
			{Title: "Name", Width: 3},
			{Title: "Type", Width: 2},
			{Title: "Null", Width: 1},
		}),
		tbl.WithRows(lo.Map(t.Columns, func(col inspect.Column, _ int) tbl.Row {
			return tbl.Row{
				formatColumnName(t, col),
				col.Type,
				formatBool(col.Null),
			}
		})),
		tbl.WithFlexColumnWidth(true),
		tbl.WithFocused(true),
	)

	idxTbl := tbl.New(
		tbl.WithColumns([]tbl.Column{
			{Title: "Name", Width: 5},
			{Title: "Unique", Width: 1},
			{Title: "Parts", Width: 3},
		}),
		tbl.WithRows(lo.Map(t.Indexes, func(idx inspect.Index, _ int) tbl.Row {
			return tbl.Row{
				idx.Name,
				formatBool(idx.Unique),
				strings.Join(lo.Map(idx.Parts, func(part inspect.IndexPart, _ int) string {
					return part.Column
				}), inlineListSeparator),
			}
		})),
		tbl.WithFlexColumnWidth(true),
		tbl.WithFocused(true),
	)

	fksTbl := tbl.New(
		tbl.WithColumns([]tbl.Column{
			{Title: "Name", Width: 2},
			{Title: "Columns", Width: 1},
			{Title: "References", Width: 1},
		}),
		tbl.WithRows(lo.Map(t.ForeignKeys, func(fk inspect.ForeignKey, _ int) tbl.Row {
			return tbl.Row{
				fk.Name,
				strings.Join(fk.Columns, inlineListSeparator),
				fmt.Sprintf("%s(%s)", fk.References.Table, strings.Join(fk.References.Columns, inlineListSeparator)),
			}
		})),
		tbl.WithFlexColumnWidth(true),
		tbl.WithFocused(true),
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
