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
  describe Describes a job
  list     Lists all jobs available
  results  Gets the results of a job's build
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	var cmd *command.Command

	if len(os.Args) == 1 {
		cmd.UsageAndExit("jk: You need to provide a command.\n")
	}

	switch os.Args[1] {
	case "build":
		cmd = command.NewBuildCommand()
	case "describe":
		cmd = command.NewDescribeCommand()
	case "list":
		cmd = command.NewListCommand()
	case "results":
		cmd = command.NewResultsCommand()
	default:
		cmd.UsageAndExit(fmt.Sprintf("jk: '%s' is not a jk comand.\n", os.Args[1]))
	}

	cmd.Init(os.Args[2:])
	cmd.Run()
}
