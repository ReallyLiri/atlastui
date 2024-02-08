package tui

import "github.com/charmbracelet/bubbles/list"

type focusedComponent int

const (
	tablesListFocused focusedComponent = iota
	detailsTabFocused
	detailsContentsFocused
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
