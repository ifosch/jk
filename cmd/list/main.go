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
