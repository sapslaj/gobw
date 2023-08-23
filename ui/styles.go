package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	logo          = ""
	selectedColor = lipgloss.Color("4")
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	mutedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("4"))
	docStyle      = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle   = focusedStyle.Copy()
	noStyle       = lipgloss.NewStyle()

	rowStyle = lipgloss.NewStyle().
			Padding(0, 0, 0, 2)

	selectedRowStyle = selectedStyle.Copy().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("4")).
				Foreground(lipgloss.Color("4")).
				Padding(0, 0, 0, 1)

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)
