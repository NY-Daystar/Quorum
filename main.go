package main

import (
	"fmt"
	"os"
	"quorum/console"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(console.Start())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
