package tui

import (
	"github.com/reallyliri/atlastui/inspect"
	"github.com/samber/lo"
	"strings"
)

const (
	breadcrumbsSeparator = "►"
	inlineListSeparator  = ", "
	tabsSeparator        = " · "
)

func formatBool(b bool) string {
	return lo.Ternary(b, "Yes", "No")
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
