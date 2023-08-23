package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	logo                = ""
	listSelectedStyle   = lipgloss.Color("4")
	titleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("4"))
	docStyle            = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)
