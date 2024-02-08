package tui

import (
	_ "embed"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	chart "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"strings"
)

//go:embed res/atlas.ans
var atlas string

func (m *model) View() string {
	if m.state.quitting || m.state.termWidth == 0 || m.state.termHeight == 0 {
		return ""
	}
	if m.state.easteregg {
		height := lipgloss.Height(atlas)
		if height > m.state.termHeight {
			atlas = strings.Join(strings.Split(atlas, "\n")[:m.state.termHeight], "\n")
		}
		return lipgloss.NewStyle().Width(m.state.termWidth).Height(m.state.termHeight).Align(lipgloss.Right).Background(whiteTint).Foreground(whiteTint).Render(atlas)
	}

	borderWidth, borderHeight := borderFocusedStyle.GetFrameSize()

	title := titleView(m.config.title, m.state.selectedSchema, m.state.selectedTable, m.vms.globe)
	footer := m.vms.help.View(m.config.keymap)
	centerHeight := m.state.termHeight - lipgloss.Height(title) - lipgloss.Height(footer) - 5

	var tablesList string
	var tabsView string
	var details string

	if m.state.selectedSchema != "" {
		m.vms.tablesList.SetSize(m.state.termWidth/3-borderWidth, centerHeight-borderHeight+3)
		tablesList = withBorder(m.vms.tablesList.View(), m.state.focused == tablesListFocused)

		if m.state.selectedTable != "" {
			detailsWidth := (m.state.termWidth*2)/3 - borderWidth
			tabsView = m.tabsView(detailsWidth, m.state.focused == detailsTabFocused)

			var currTbl chart.Model
			switch m.state.selectedTab {
			case ColumnsTable:
				currTbl = m.vms.colsChart
			case IndexesTable:
				currTbl = m.vms.idxChart
			case ForeignKeysTable:
				currTbl = m.vms.fksChart
			}
			currTbl.SetWidth(detailsWidth)
			currTbl.SetHeight(centerHeight - lipgloss.Height(tabsView) - borderHeight + 2)
			if len(currTbl.Rows()) == 0 {
				noData := fmt.Sprintf("No %s", m.state.selectedTab.title())
				details = withBorder(noDataStyle.Copy().
					Width(currTbl.Width()).
					Height(currTbl.Height()+1).
					Render(noData), m.state.focused == detailsContentsFocused)
			} else {
				details = withBorder(currTbl.View(), m.state.focused == detailsContentsFocused)
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
		tabView(ColumnsTable.title(), m.state.selectedTab == ColumnsTable),
		tabView(IndexesTable.title(), m.state.selectedTab == IndexesTable),
		tabView(ForeignKeysTable.title(), m.state.selectedTab == ForeignKeysTable),
	}

	row := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(strings.Join(tabs, subTitleStyle.Render(tabsSeparator)))
	return withBorder(row, focused)
}

func tabView(title string, selected bool) string {
	return lo.Ternary(selected, titleStyle.Render(title), subTitleStyle.Render(title))
}

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

func withBorder(ui string, focused bool) string {
	if focused {
		return borderFocusedStyle.Render(ui)
	} else {
		return borderBluredStyle.Render(ui)
	}
}
