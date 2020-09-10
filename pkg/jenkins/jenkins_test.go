package jenkins

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/bndr/gojenkins"
)

func TestJenkins_NewJenkins(t *testing.T) {
	tcs := []struct {
		url           string
		expectedError error
	}{
		{
			"",
			fmt.Errorf("Get \"/api/json\": unsupported protocol scheme \"\""),
		},
		{
			"http://google.com",
			fmt.Errorf("Jenkins server version is empty"),
		},
	}
	for _, tc := range tcs {
		_, err := NewJenkins(tc.url, "admin", "admin", nil)
		if err != nil {
			if tc.expectedError != nil {
				if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tc.expectedError) {
					t.Fatalf("got %v error, but expected %v", err, tc.expectedError)
				}
			} else {
				t.Fatal(err)
			}
		}
	}
}

func TestJenkins_List(t *testing.T) {
	tcs := []struct {
		jobNameList []string
		messages    []Message
	}{
		{
			[]string{
				"job1",
				"job2",
			},
			[]Message{
				{
					Message: "job1",
					Error:   false,
					Done:    false,
				},
				{
					Message: "job2",
					Error:   false,
					Done:    false,
				},
				{
					Message: "",
					Error:   false,
					Done:    true,
				},
			},
		},
	}

	j := &Jenkins{client: jenkinsClientMock{}}
	for _, tc := range tcs {
		getAllJobsMock = func() (jobs []*gojenkins.Job, err error) {
			for _, jobName := range tc.jobNameList {
				jobs = append(jobs, &gojenkins.Job{
					Raw: &gojenkins.JobResponse{
						Name: jobName,
					},
				},
				)
			}
			return jobs, nil
		}

		channel := make(chan Message)
		defer close(channel)
		go j.List(channel)

		assertExpectedMessages(tc.messages, channel, t)
	}
}

func TestJenkins_Describe(t *testing.T) {
	tcs := []struct {
		jobName  string
		tmpl     string
		messages []Message
	}{
		{
			"job1",
			"{{ .Raw.Name }}",
			[]Message{
				{
					Message: "job1",
					Error:   false,
					Done:    true,
				},
			},
		},
	}

	j := &Jenkins{client: jenkinsClientMock{}}
	for _, tc := range tcs {
		getJobMock = func(jobName string, parents ...string) (job *gojenkins.Job, err error) {
			return &gojenkins.Job{
				Raw: &gojenkins.JobResponse{
					Name: tc.jobName,
				},
			}, nil
		}
		templ, err := template.New("Job").Parse(tc.tmpl)
		if err != nil {
			panic(fmt.Sprintf("Template Creation error %v", err))
		}

		channel := make(chan Message)
		defer close(channel)
		go j.Describe(tc.jobName, templ, channel)

		assertExpectedMessages(tc.messages, channel, t)
	}
}

func TestJenkins_Build(t *testing.T) {
	tcs := []struct {
		jobName  string
		params   map[string]string
		item     int64
		buildID  int64
		messages []Message
	}{
		{
			"job_name",
			map[string]string{},
			1000,
			100,
			[]Message{
				{
					Message: "Build queued /job/job_name/100",
					Error:   false,
					Done:    false,
				},
				{
					Message: "Build started /job/job_name/100",
					Error:   false,
					Done:    false,
				},
				{
					Message: "Build finished /job/job_name/100",
					Error:   false,
					Done:    true,
				},
			},
		},
		{
			"job_with_params",
			map[string]string{
				"ARG1": "value1",
			},
			2000,
			200,
			[]Message{
				{
					Message: "Build queued /job/job_with_params/200",
					Error:   false,
					Done:    false,
				},
				{
					Message: "Build started /job/job_with_params/200",
					Error:   false,
					Done:    false,
				},
				{
					Message: "Build finished /job/job_with_params/200",
					Error:   false,
					Done:    true,
				},
			},
		},
	}

	j := &Jenkins{client: jenkinsClientMock{}}
	for _, tc := range tcs {
		getQueueItemMock = func(number int64) (task *Task, err error) {
			return &Task{BuildID: tc.buildID}, nil
		}
		waitForBuild = func(build *gojenkins.Build) (err error) {
			build.Raw.Building = false
			return nil
		}
		getBuildMock = func(jobName string, buildID int64) (build *gojenkins.Build, err error) {
			return &gojenkins.Build{
				Base: fmt.Sprintf("/job/%v/%v", jobName, buildID),
				Raw: &gojenkins.BuildResponse{
					Building: true,
				},
			}, nil
		}
		buildJobMock = func(jobName string, option ...interface{}) (queueItem int64, err error) {
			return tc.item, nil
		}

		channel := make(chan Message)
		defer close(channel)
		go j.Build(tc.jobName, tc.params, channel)

		assertExpectedMessages(tc.messages, channel, t)
	}
}
