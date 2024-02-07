package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/samber/lo"
)

type tablesListItem string

var _ list.DefaultItem = tablesListItem("")

func (t tablesListItem) FilterValue() string {
	return string(t)
}

func (t tablesListItem) Title() string {
	return string(t)
}

func (t tablesListItem) Description() string {
	return ""
}

func newTablesList(names []string) list.Model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	lst := list.New(
		lo.Map(names, func(name string, _ int) list.Item {
			return tablesListItem(name)
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
