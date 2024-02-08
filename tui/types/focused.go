package types

type FocusedComponent int

const (
	TablesListFocused FocusedComponent = iota
	DetailsTabFocused
	DetailsContentsFocused
)
