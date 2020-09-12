package jenkins

import "github.com/bndr/gojenkins"

type jenkinsClientMock struct {
	jobs []*gojenkins.Job
}

func newJenkinsClientMock(jobNames []string) (j jenkinsClientMock) {
	j = jenkinsClientMock{}
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

var getBuildMock func(jobŃame string, buildId int64) (build *gojenkins.Build, err error)

func (j jenkinsClientMock) GetBuild(jobName string, buildID int64) (build *gojenkins.Build, err error) {
	return getBuildMock(jobName, buildID)
}

var buildJobMock func(jobName string, option ...interface{}) (queueItem int64, err error)

func (j jenkinsClientMock) BuildJob(jobName string, options ...interface{}) (queueItem int64, err error) {
	return buildJobMock(jobName, options)
}

func (j jenkinsClientMock) GetAllJobs() (jobs []*gojenkins.Job, err error) {
	return j.jobs, nil
}

var getJobMock func(jobName string, parents ...string) (job *gojenkins.Job, err error)

func (j jenkinsClientMock) GetJob(jobName string, parents ...string) (job *gojenkins.Job, err error) {
	return getJobMock(jobName, parents...)
}

var getQueueItemMock func(int64) (*Task, error)

func (j jenkinsClientMock) GetQueueItem(number int64) (task *Task, err error) {
	return getQueueItemMock(number)
}
