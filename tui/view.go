package tui

import (
	_ "embed"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	chart "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/reallyliri/atlastui/tui/format"
	"github.com/reallyliri/atlastui/tui/styles"
	"github.com/reallyliri/atlastui/tui/types"
	"github.com/samber/lo"
	"strings"
)

func (m *model) View() string {
	if m.state.quitting || m.state.termWidth == 0 || m.state.termHeight == 0 {
		return ""
	}

	borderWidth, borderHeight := styles.BorderFocusedStyle.GetFrameSize()

	title := titleView(m.config.title, m.state.selectedSchema, m.state.selectedTable, m.vms.globe)
	footer := m.vms.help.View(m.config.keymap)
	centerHeight := m.state.termHeight - lipgloss.Height(title) - lipgloss.Height(footer) - 5

	var tablesList string
	var tabsView string
	var details string

	if m.state.selectedSchema != "" {
		m.vms.tablesList.SetSize(m.state.termWidth/3-borderWidth, centerHeight-borderHeight+3)
		tablesList = withBorder(m.vms.tablesList.View(), m.state.focused == types.TablesListFocused)

		if m.state.selectedTable != "" {
			detailsWidth := (m.state.termWidth*2)/3 - borderWidth
			tabsView = m.tabsView(detailsWidth, m.state.focused == types.DetailsTabFocused)

			var currChart chart.Model
			switch m.state.selectedTab {
			case types.ColumnsTable:
				currChart = m.vms.colsChart
			case types.IndexesTable:
				currChart = m.vms.idxChart
			case types.ForeignKeysTable:
				currChart = m.vms.fksChart
			}
			currChart.SetWidth(detailsWidth)
			currChart.SetHeight(centerHeight - lipgloss.Height(tabsView) - borderHeight + 2)
			if len(currChart.Rows()) == 0 {
				noData := fmt.Sprintf("No %s", m.state.selectedTab.Title())
				details = withBorder(styles.NoDataStyle.Copy().
					Width(currChart.Width()).
					Height(currChart.Height()+1).
					Render(noData), m.state.focused == types.DetailsContentsFocused)
			} else {
				details = withBorder(currChart.View(), m.state.focused == types.DetailsContentsFocused)
			}
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			tablesList,
			lipgloss.JoinVertical(
				lipgloss.Top,
				tabsView, details,
			),
		),
		footer,
	)
}

func (m *model) tabsView(width int, focused bool) string {
	tabs := []string{
		tabView(types.ColumnsTable.Title(), m.state.selectedTab == types.ColumnsTable),
		tabView(types.IndexesTable.Title(), m.state.selectedTab == types.IndexesTable),
		tabView(types.ForeignKeysTable.Title(), m.state.selectedTab == types.ForeignKeysTable),
	}

	row := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(strings.Join(tabs, styles.SubTitleStyle.Render(format.TabsSeparator)))
	return withBorder(row, focused)
}

func tabView(title string, selected bool) string {
	return lo.Ternary(selected, styles.TitleStyle.Render(title), styles.SubTitleStyle.Render(title))
}

func titleView(title, selectedSchema, selectedTable string, globe spinner.Model) string {
	parts := make([]string, 0, 8)
	parts = append(parts, styles.TitleStyle.Render(title))
	parts = append(parts, " ")
	parts = append(parts, globe.View())
	if selectedSchema != "" {
		parts = append(parts, styles.BreadcrumbsSectionStyle.Render(" ", format.BreadcrumbsSeparator, " schema "))
		parts = append(parts, styles.BreadcrumbsTitleStyle.Render(selectedSchema))
	}
	if selectedTable != "" {
		parts = append(parts, styles.BreadcrumbsSectionStyle.Render(" ", format.BreadcrumbsSeparator, " table "))
		parts = append(parts, styles.BreadcrumbsTitleStyle.Render(selectedTable))
	}
	parts = append(parts, styles.BreadcrumbsSectionStyle.Render(" ", format.BreadcrumbsSeparator))
	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}

func withBorder(ui string, focused bool) string {
	if focused {
		return styles.BorderFocusedStyle.Render(ui)
	} else {
		return styles.BorderBluredStyle.Render(ui)
	}
}
