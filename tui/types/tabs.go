package types

type TableDetailsSection int

const (
	ColumnsTable TableDetailsSection = iota
	IndexesTable
	ForeignKeysTable
)

func (section TableDetailsSection) Title() string {
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
