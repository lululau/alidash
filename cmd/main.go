package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"aliyun-tui-viewer/internal/tui"
)

func main() {
	// Create new application model
	model, err := tui.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Create and run the Bubble Tea program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
