package jenkins

import (
	"errors"
	"fmt"

	"github.com/bndr/gojenkins"
)

type jenkinsClientMock struct {
	jobs        []*gojenkins.Job
	nextBuildID int64
	nextItemID  int64
	Server      string
	Version     string
}

func newJenkinsClientMock(jobNames []string, nextBuildID int64, nextItemID int64) (j *jenkinsClientMock) {
	j = &jenkinsClientMock{
		nextBuildID: nextBuildID,
		nextItemID:  nextItemID,
		Server:      "http://mockedjenkins",
		Version:     "",
	}
	for _, name := range jobNames {
		j.jobs = append(
			j.jobs,
			&gojenkins.Job{
				Raw: &gojenkins.JobResponse{
					Name: name,
				},
			},
		)
	}
	return
}

func (j jenkinsClientMock) GetBuild(jobName string, buildID int64) (build *Build, err error) {
	return NewBuild(
		&gojenkins.Build{
			Base: fmt.Sprintf("/job/%v/%v", jobName, buildID),
			Jenkins: &gojenkins.Jenkins{
				Server:  j.Server,
				Version: j.Version,
			},
			Job: &gojenkins.Job{
				Raw: &gojenkins.JobResponse{
					Name: jobName,
				},
			},
			Raw: &gojenkins.BuildResponse{
				Building: true,
				ID:       fmt.Sprintf("%v", buildID),
			},
		},
	), nil
}

func (j jenkinsClientMock) BuildJob(jobName string, options ...interface{}) (queueItem int64, err error) {
	return j.nextItemID, nil
}

func (j jenkinsClientMock) GetAllJobs() (jobs []*gojenkins.Job, err error) {
	return j.jobs, nil
}

func (j jenkinsClientMock) GetJob(jobName string, parents ...string) (foundJob *gojenkins.Job, err error) {
	for _, job := range j.jobs {
		if jobName == job.GetName() {
			return job, nil
		}
	}
	return nil, errors.New("404")
}

func (j *jenkinsClientMock) GetQueueItem(number int64) (task *Task, err error) {
	task = &Task{BuildID: j.nextBuildID}
	j.nextBuildID++
	return
}

func (j jenkinsClientMock) GetLastBuild(string) (*Build, error) {
	return nil, nil
}
