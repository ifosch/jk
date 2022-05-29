package command

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ifosch/jk/pkg/jenkins"
)

var describeUsage = `Describes a job.

Usage: jk describe <jobName>
`

// NewDescribeCommand returns the command for the Describe operation.
func NewDescribeCommand() *Command {
	cmd := &Command{
		flags:   flag.NewFlagSet("describe", flag.ExitOnError),
		Execute: describeFunc,
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, describeUsage)
	}

	return cmd
}

var describeFunc = func(cmd *Command, args []string) {
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

	if len(args) < 1 {
		log.Fatal("Missing job name and optional params")
	}
	jobName := args[0]
	go j.Describe(jobName, nil, ch)

	var output string
	for message := range ch {
		output = message.Message
		if message.Done {
			break
		}
	}
	fmt.Println(output)
}
