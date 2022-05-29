package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ifosch/jk/pkg/jenkins"
)

func main() {
	url, ok := os.LookupEnv("JENKINS_URL")
	if !ok {
		url = "http://localhost:8080"
	}
	user, ok := os.LookupEnv("JENKINS_USER")
	if !ok {
		user = "admin"
	}
	password, ok := os.LookupEnv("JENKINS_PASSWORD")
	if !ok {
		password = "admin"
	}

	j, err := jenkins.NewJenkins(url, user, password, nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("You need to specify at least the job name")
	}

	jobName := os.Args[1]
	var buildID int64 = 0
	if len(os.Args) > 2 {
		buildID, err = strconv.ParseInt(os.Args[2], 10, 64)
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
