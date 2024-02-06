package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

const (
	blueTint  = lipgloss.Color("#388de9")
	greenTint = lipgloss.Color("#46b17b")
)

const (
	breadcrumbsSeparator = "►"
	inlineListSeparator  = ", "
)

var (
	titleStyle              = lipgloss.NewStyle().Foreground(blueTint)
	breadcrumbsSectionStyle = help.New().Styles.ShortKey
	breadcrumbsTitleStyle   = lipgloss.NewStyle().Foreground(greenTint)
)

func formatBool(b bool) string {
	return lo.Ternary(b, "Yes", "No")
}
