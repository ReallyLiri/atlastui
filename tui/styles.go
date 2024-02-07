package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

const (
	whiteTint         = lipgloss.Color("#f0f0f0")
	blueTint          = lipgloss.Color("#388de9")
	greenTint         = lipgloss.Color("#46b17b")
	borderFocusedTint = lipgloss.Color("63")
	borderBluredTint  = lipgloss.Color("240")
)

const (
	breadcrumbsSeparator = "►"
	inlineListSeparator  = ", "
)

var (
	subTitleTint = help.New().Styles.ShortKey.GetForeground()

	titleStyle              = lipgloss.NewStyle().Foreground(blueTint)
	subTitleStyle           = lipgloss.NewStyle().Foreground(subTitleTint)
	breadcrumbsSectionStyle = subTitleStyle
	breadcrumbsTitleStyle   = lipgloss.NewStyle().Foreground(greenTint)
	borderFocusedStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(borderFocusedTint)
	borderBluredStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(borderBluredTint)
	alignLeftStyle          = lipgloss.NewStyle().Align(lipgloss.Left)
	boldStyle               = lipgloss.NewStyle().Bold(true)
)

func formatBool(b bool) string {
	return lo.Ternary(b, "Yes", "No")
}
