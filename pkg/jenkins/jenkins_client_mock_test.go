package jenkins

import "github.com/bndr/gojenkins"

type jenkinsClientMock struct{}

var getBuildMock func(jobŃame string, buildId int64) (build *gojenkins.Build, err error)

func (j jenkinsClientMock) GetBuild(jobName string, buildID int64) (build *gojenkins.Build, err error) {
	return getBuildMock(jobName, buildID)
}

var buildJobMock func(jobName string, option ...interface{}) (queueItem int64, err error)

func (j jenkinsClientMock) BuildJob(jobName string, options ...interface{}) (queueItem int64, err error) {
	return buildJobMock(jobName, options)
}

var getAllJobsMock func() (jobs []*gojenkins.Job, err error)

func (j jenkinsClientMock) GetAllJobs() (jobs []*gojenkins.Job, err error) {
	return getAllJobsMock()
}

var getJobMock func(jobName string, parents ...string) (job *gojenkins.Job, err error)

func (j jenkinsClientMock) GetJob(jobName string, parents ...string) (job *gojenkins.Job, err error) {
	return getJobMock(jobName, parents...)
}
