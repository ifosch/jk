package jenkins

import "github.com/bndr/gojenkins"

// API is an interface defining a Jenkins API client methods.
type API interface {
	GetBuild(string, int64) (*Build, error)
	BuildJob(string, ...interface{}) (int64, error)
	GetAllJobs() ([]*gojenkins.Job, error)
	GetJob(string, ...string) (*gojenkins.Job, error)
	GetQueueItem(int64) (*Task, error)
	GetLastBuild(string) (*Build, error)
}
