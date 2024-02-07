package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func titleView(title, selectedSchema, selectedTable string, globe spinner.Model) string {
	parts := make([]string, 0, 8)
	parts = append(parts, titleStyle.Render(title))
	parts = append(parts, " ")
	parts = append(parts, globe.View())
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
