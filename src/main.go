package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	buildVersion = "dev"
	buildCommit  = "dirty"
	argv0        = "paisano"
	project      = "Paisano"
)

func main() {
	if len(os.Args[1:]) == 0 {
		// with NO arguments, invoke the TUI
		if model, err := tea.NewProgram(
			InitialPage(),
			tea.WithAltScreen(),
		).StartReturningModel(); err != nil {
			log.Fatalf("Error running program: %s", err)
		} else if err := model.(*Tui).FatalError; err != nil {
			log.Fatal(err)
		} else if command := model.(*Tui).ExecveCommand; command != nil {
			if err := command.Exec(nil); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		// with arguments, invoke the CLI
		ExecuteCli()
	}
}
