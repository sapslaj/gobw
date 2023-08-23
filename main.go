package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
)

var bwm BWManager

func main() {
	cmd := exec.Command("bw", "-v")
	if err := cmd.Run(); err != nil {
		fmt.Println("Could not find 'bw' command in '$PATH'. Please check if Bitwarden CLI is installed.\nGoodbye")
		os.Exit(1)
	}
	if clipboard.Unsupported {
		// TODO: better error message
		panic("failed to setup clipboard.")
	}
	bwm = NewBWManager()
	m := NewMainModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
