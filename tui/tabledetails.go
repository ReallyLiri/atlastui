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

func newTblDetails(t inspect.Table, width, height int) (tbl.Model, tbl.Model, tbl.Model) {
	colsTbl := tbl.New(
		tbl.WithColumns([]tbl.Column{
			{Title: "Name", Width: width / 2},
			{Title: "Type", Width: width / 3},
			{Title: "Null", Width: width / 6},
		}),
		tbl.WithRows(lo.Map(t.Columns, func(col inspect.Column, _ int) tbl.Row {
			return tbl.Row{
				col.Name,
				col.Type,
				formatBool(col.Null),
			}
		})),
		tbl.WithFocused(true),
		tbl.WithWidth(width),
		tbl.WithHeight(height),
	)

	idxTbl := tbl.New(
		tbl.WithColumns([]tbl.Column{
			{Title: "Name", Width: 5 * width / 9},
			{Title: "Unique", Width: width / 9},
			{Title: "Parts", Width: 3 * width / 9},
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
		tbl.WithWidth(width),
		tbl.WithHeight(height),
	)

	fksTbl := tbl.New(
		tbl.WithColumns([]tbl.Column{
			{Title: "Name", Width: width / 2},
			{Title: "Columns", Width: width / 4},
			{Title: "References", Width: width / 4},
		}),
		tbl.WithRows(lo.Map(t.ForeignKeys, func(fk inspect.ForeignKey, _ int) tbl.Row {
			return tbl.Row{
				fk.Name,
				strings.Join(fk.Columns, inlineListSeparator),
				fmt.Sprintf("%s(%s)", fk.References.Table, strings.Join(fk.References.Columns, inlineListSeparator)),
			}
		})),
		tbl.WithWidth(width),
		tbl.WithHeight(height),
	)

	return colsTbl, idxTbl, fksTbl
}
