package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ifosch/jk/pkg/command"
)

var usage = `Usage: jk command [options]

A jk CLI focused in job management tasks.

Commands:
  build    Builds a job
  list     Lists all jobs available
`

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprint(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	var cmd *command.Command

	if len(os.Args) == 1 {
		usageAndExit("jk: You need to provide a command.\n")
	}

	switch os.Args[1] {
	case "build":
		cmd = command.NewBuildCommand()
	case "list":
		cmd = command.NewListCommand()
	default:
		usageAndExit(fmt.Sprintf("jk: '%s' is not a jk comand.\n", os.Args[1]))
	}

	cmd.Init(os.Args[2:])
	cmd.Run()
}
