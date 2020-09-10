package jenkins

import "github.com/bndr/gojenkins"

// Task ...
type Task struct {
	JenkinsTask *gojenkins.Task
	BuildID     int64
}

// Poll ...
func (t *Task) Poll() {
	t.JenkinsTask.Poll()
	t.BuildID = t.JenkinsTask.Raw.Executable.Number
}
