package command

import (
	"flag"
	"fmt"
	"os"
)

// Command represents a command for mecha.
type Command struct {
	flags   *flag.FlagSet
	Execute func(cmd *Command, args []string)
}

// Init method receives the args for the command and parses them into
// options.
func (c *Command) Init(args []string) error {
	return c.flags.Parse(args)
}

// Run executes the command code.
func (c *Command) Run() {
	c.Execute(c, c.flags.Args())
}

func (c *Command) usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprint(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n")
	}

	c.flags.Usage()
	os.Exit(1)
}
