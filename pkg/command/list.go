package command

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ifosch/jk/pkg/jenkins"
)

var listUsage = `Lists jobs.

Usage: jk list
`

// NewListCommand returns the command for the List operation.
func NewListCommand() *Command {
	cmd := &Command{
		flags:   flag.NewFlagSet("list", flag.ExitOnError),
		Execute: listFunc,
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, listUsage)
	}

	return cmd
}

var listFunc = func(cmd *Command, args []string) {
	j, err := jenkins.NewJenkins(
		os.Getenv("JENKINS_URL"),
		os.Getenv("JENKINS_USER"),
		os.Getenv("JENKINS_PASSWORD"),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan jenkins.Message)
	defer close(ch)

	go j.List(ch)

	for message := range ch {
		if !message.Done {
			fmt.Println(message.Message)
		}
		if message.Done {
			break
		}
	}
}
