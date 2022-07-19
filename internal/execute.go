package internal

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Execute() {
	p := tea.NewProgram(NewModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
