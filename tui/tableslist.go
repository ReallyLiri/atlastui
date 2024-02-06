package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/samber/lo"
)

func newTablesList(names []string, width, height int) table.Model {
	return table.New(
		table.WithColumns([]table.Column{
			{Title: "Name", Width: width},
		}),
		table.WithRows(lo.Map(names, func(name string, _ int) table.Row {
			return table.Row{
				name,
			}
		})),
		table.WithWidth(width),
		table.WithHeight(height),
	)
}
