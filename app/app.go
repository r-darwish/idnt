package app

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func Run() {
	m := newMainScreen()

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
