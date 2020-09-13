package jenkins

import (
	"fmt"
	"net/http"

	"github.com/bndr/gojenkins"
)

// Client ...
type Client struct {
	*gojenkins.Jenkins
}

// NewClient ...
func NewClient(url, user, password string, client *http.Client) (c *Client, err error) {
	j, err := gojenkins.CreateJenkins(
		client,
		url,
		user,
		password,
	).Init()
	if err != nil {
		return
	}
	if j.Version == "" {
		return nil, fmt.Errorf("Jenkins server version is empty")
	}
	return &Client{
		j,
	}, nil
}

// GetQueueItem ...
func (c *Client) GetQueueItem(number int64) (task *Task, err error) {
	jenkinsTask, err := c.Jenkins.GetQueueItem(number)
	if err != nil {
		return nil, err
	}
	return &Task{
		JenkinsTask: jenkinsTask,
		BuildID:     jenkinsTask.Raw.Executable.Number,
	}, nil
}

// GetBuild ...
func (c *Client) GetBuild(jobName string, buildID int64) (build *Build, err error) {
	jenkinsBuild, err := c.Jenkins.GetBuild(jobName, buildID)
	if err != nil {
		return
	}
	build = NewBuild(jenkinsBuild)
	return
}
