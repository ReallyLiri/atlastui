package tui

import "github.com/charmbracelet/lipgloss"

func titleView(title, selectedSchema, selectedTable string) string {
	parts := make([]string, 0, 7)
	parts = append(parts, titleStyle.Render(title))
	parts = append(parts, " 🌍")
	if selectedSchema != "" {
		parts = append(parts, breadcrumbsSectionStyle.Render(" ", breadcrumbsSeparator, " schema "))
		parts = append(parts, breadcrumbsTitleStyle.Render(selectedSchema))
	}
	if selectedTable != "" {
		parts = append(parts, breadcrumbsSectionStyle.Render(" ", breadcrumbsSeparator, " table "))
		parts = append(parts, breadcrumbsTitleStyle.Render(selectedTable))
	}
	parts = append(parts, breadcrumbsSectionStyle.Render(" ", breadcrumbsSeparator))
	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}
