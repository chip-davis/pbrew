package main

import (
	"fmt"
	"os"

	"github.com/chipdavis/pbrew/internal/brew"
	"github.com/chipdavis/pbrew/internal/commands"
	"github.com/chipdavis/pbrew/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "list", "ls":
		err = commands.List(args)
	case "outdated":
		err = commands.Outdated(args)
	case "search":
		err = commands.Search(args)
	case "info":
		err = commands.Info(args)
	case "update":
		err = commands.Update(args)
	case "upgrade":
		err = commands.Upgrade(args)
	case "help", "-h", "--help":
		printUsage()
	default:
		// Unknown commands fall through to brew so things like `pbrew install foo` still work.
		err = brew.Passthrough(append([]string{cmd}, args...))
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render(err.Error()))
		os.Exit(1)
	}
}

func printUsage() {
	// TODO: implement pretty usage output using ui styles
}
