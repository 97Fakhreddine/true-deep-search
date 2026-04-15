package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"hybridsearch/internal/app"
)

func main() {
	model := app.NewModel()

	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		log.Fatalf("failed to run hybridsearch: %v", err)
	}
}
