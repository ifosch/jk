package jenkins

import (
	"fmt"
	"os"
)

// Connect establishes a connection to a Jenkins.
func Connect() (j *Jenkins, err error) {
	var url, user, password string
	var ok bool
	if url, ok = os.LookupEnv("JENKINS_URL"); !ok {
		url = "http://localhost:8080"
	}
	if user, ok = os.LookupEnv("JENKINS_USER"); !ok {
		user = "admin"
	}
	if password, ok = os.LookupEnv("JENKINS_PASSWORD"); !ok {
		password = "admin"
	}
	j, err = NewJenkins(url, user, password, nil)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Jenkins at %v as %v:\n%v", url, user, err)
	}
	return j, nil
}
