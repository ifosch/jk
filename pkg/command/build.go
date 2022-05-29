package command

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ifosch/jk/pkg/jenkins"
)

var buildUsage = `Builds a job.

Usage: jk build <jobName> [<params>...]
`

// NewBuildCommand returns the command for the Build operation.
func NewBuildCommand() *Command {
	cmd := &Command{
		flags:   flag.NewFlagSet("build", flag.ExitOnError),
		Execute: buildFunc,
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, buildUsage)
	}

	return cmd
}

var buildFunc = func(cmd *Command, args []string) {
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
	params := map[string]string{}
	for _, arg := range args[1:] {
		keyValue := strings.Split(arg, "=")
		params[keyValue[0]] = keyValue[1]
	}
	go j.Build(jobName, params, ch)

	for message := range ch {
		fmt.Println(message.Message)
		if message.Done {
			break
		}
	}
}
