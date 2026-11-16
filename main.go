package main

import (
	"fmt"
	"statboard/cmd"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		Render("StatBoard")
	fmt.Println(title)
	cmd.Execute()
}
