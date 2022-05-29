// Package jenkins provides a nice and easy to use interface with
// Jenkins job execution platform.
package jenkins

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/ifosch/jk/pkg/templates"
)

// Jenkins is the Jenkins object.
type Jenkins struct {
	client API
}

// NewJenkins initializes a Jenkins object with the corresponding url,
// credentials, and HTTP client, if provided. The HTTP client can be
// nil, and a default one will be created by the API. Returns a
// reference to the Jenkins object, and nil, if no error happened, or
// nil, and the error, otherwise.
func NewJenkins(
	url, user, password string,
	client *http.Client,
) (jenkins *Jenkins, err error) {
	c, err := NewClient(
		url,
		user,
		password,
		client,
	)
	if err != nil {
		return
	}
	return &Jenkins{
		client: c,
	}, nil
}

// List all jobs. It will use the provided channel to report one
// message for each job. Once all jobs were sent, an empty string
// message will be issued with Done set to true.
func (j *Jenkins) List(out chan Message) {
	jobs, err := j.client.GetAllJobs()
	if err != nil {
		reply(
			fmt.Sprintf("Jenkins GetAllJobs error %v", err),
			true,
			true,
			out,
		)
	}
	for _, job := range jobs {
		reply(
			job.GetName(),
			false,
			false,
			out,
		)
	}
	reply(
		"",
		false,
		true,
		out,
	)
}

// Describe provides a templated description of a job, identified by
// jobName.
func (j *Jenkins) Describe(jobName string, t *template.Template, out chan Message) {
	if t == nil {
		var err error
		t, err = templates.Describe()
		if err != nil {
			reply(
				fmt.Sprintf("Template parsing error %v", err),
				true,
				true,
				out,
			)
		}
	}
	job, err := j.client.GetJob(jobName)
	if err != nil {
		reply(
			fmt.Sprintf("Job get error %v", err),
			true,
			true,
			out,
		)
	}
	msg := &bytes.Buffer{}
	err = t.Execute(msg, job)
	if err != nil {
		reply(
			fmt.Sprintf("Template Execute error %v", err),
			true,
			true,
			out,
		)
	}
	reply(
		msg.String(),
		false,
		true,
		out,
	)
}

// Results returns results for a specific job's build.
func (j *Jenkins) Results(jobName string, buildID int64) (r []*Case, err error) {
	var b *Build
	if buildID > 0 {
		b, err = j.client.GetBuild(jobName, buildID)
	} else {
		b, err = j.client.GetLastBuild(jobName)
	}
	if err != nil {
		return nil, err
	}

	results, err := b.GetResults()
	if err != nil {
		return nil, err
	}

	if len(results.Suites) == 0 {
		return nil, fmt.Errorf("last build (%v) for job %v not completed", b.GetBuildNumber(), jobName)
	}

	r = []*Case{}
	for _, suite := range results.Suites {
		for _, testCase := range suite.Cases {
			r = append(r, &Case{
				Name:   testCase.Name,
				Status: testCase.Status,
			})
		}
	}
	return
}

// Build executes jobName with params parameters. It will use the
// channel to reply with updates on the progress of the build.
func (j *Jenkins) Build(jobName string, params map[string]string, out chan Message) {
	number, err := j.client.BuildJob(jobName, params)
	if err != nil {
		reply(
			fmt.Sprintf("Job Invoke error %v", err),
			true,
			true,
			out,
		)
	}
	task, err := j.client.GetQueueItem(number)
	if err != nil {
		reply(
			fmt.Sprintf("Task get error %v", err),
			true,
			true,
			out,
		)
	}
	reply(
		fmt.Sprintf("Build queued for job %v", jobName),
		false,
		false,
		out,
	)
	buildID := task.BuildID
	for {
		if buildID != 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
		task.Poll()
		buildID = task.BuildID
	}
	build, err := j.client.GetBuild(jobName, buildID)
	if err != nil {
		reply(
			fmt.Sprintf("Queue item get error %v", err),
			true,
			true,
			out,
		)
	}
	reply(
		fmt.Sprintf("Build started %v", build.GetURL()),
		false,
		false,
		out,
	)
	err = build.Wait()
	if err != nil {
		reply(
			fmt.Sprintf("Error polling build %v", err),
			true,
			true,
			out,
		)
	}
	reply(
		fmt.Sprintf("Build finished %v", build.GetURL()),
		false,
		true,
		out,
	)
}
