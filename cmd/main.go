package main

import (
	"flag"

	"github.com/cycle-labs/test-app/internal/application"
)

func main() {
	isTerminal := false
	flag.BoolVar(&isTerminal, "local", false, "start local terminal application")
	flag.Parse()

	if isTerminal {
		application.StartLocalTerminalApplication()
	} else {
		application.StartServer()
	}
}
