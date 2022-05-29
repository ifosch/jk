package command

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ifosch/jk/pkg/jenkins"
)

var resultsUsage = `Returns results for a job build.

Usage: jk results <jobName> [<build no>]
`

// NewResultsCommand returns the command for the Results operation.
func NewResultsCommand() *Command {
	cmd := &Command{
		flags:   flag.NewFlagSet("results", flag.ExitOnError),
		Execute: resultsFunc,
	}

	cmd.flags.Usage = func() {
		fmt.Fprintln(os.Stderr, resultsUsage)
	}

	return cmd
}

var resultsFunc = func(cmd *Command, args []string) {
	j, err := jenkins.NewJenkins(
		os.Getenv("JENKINS_URL"),
		os.Getenv("JENKINS_USER"),
		os.Getenv("JENKINS_PASSWORD"),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) < 1 {
		log.Fatalf("You need to specify at least the job name")
	}

	jobName := args[0]
	var buildID int64 = 0
	if len(args) > 1 {
		buildID, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}

	results, err := j.Results(jobName, buildID)
	if err != nil {
		log.Fatal(err)
	}

	for _, testCase := range results {
		fmt.Println(testCase.Status, testCase.Name)
	}
}
