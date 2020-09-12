package main

import (
	"fmt"
	"log"
	"os"

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

	ch := make(chan jenkins.Message)
	defer close(ch)

	if len(os.Args) < 2 {
		log.Fatal("Missing job name and optional params")
	}
	jobName := os.Args[1]
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
