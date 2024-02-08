package types

import "github.com/charmbracelet/bubbles/list"

type TablesListItem string

var _ list.DefaultItem = TablesListItem("")

func (t TablesListItem) FilterValue() string {
	return string(t)
}

func (t TablesListItem) Title() string {
	return string(t)
}

func (t TablesListItem) Description() string {
	return ""
}
