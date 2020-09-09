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
		t.Log(tc.url)
		_, err := NewJenkins(tc.url, "admin", "admin", nil)
		if err != nil {
			if tc.expectedError != nil {
				t.Log(err)
				t.Log(tc.expectedError)
				if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tc.expectedError) {
					t.Fatalf("got %v error, but expected %v", err, tc.expectedError)
				}
			} else {
				t.Fatal(err)
			}
		}
	}
}

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

		for _, expectedMsg := range tc.messages {
			t.Log(expectedMsg.Message)
			msg := <-channel
			if msg.Message != expectedMsg.Message {
				t.Fatalf("Unexpected Message in reply: %s != %s", msg.Message, expectedMsg.Message)
			}
			if msg.Error != expectedMsg.Error {
				t.Fatalf("Unexpected Error in reply: %v != %v", msg.Error, expectedMsg.Error)
			}
			if msg.Done != expectedMsg.Done {
				t.Fatalf("Unexpected Done in reply: %v != %v", msg.Done, expectedMsg.Done)
			}
		}
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

		for _, expectedMsg := range tc.messages {
			t.Log(expectedMsg.Message)
			msg := <-channel
			if msg.Message != expectedMsg.Message {
				t.Fatalf("Unexpected Message in reply: %s != %s", msg.Message, expectedMsg.Message)
			}
			if msg.Error != expectedMsg.Error {
				t.Fatalf("Unexpected Error in reply: %v != %v", msg.Error, expectedMsg.Error)
			}
			if msg.Done != expectedMsg.Done {
				t.Fatalf("Unexpected Done in reply: %v != %v", msg.Done, expectedMsg.Done)
			}
		}
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
		getQueueItemInfo = func(j *Jenkins, queueItemURL string) (buildID int64, err error) {
			if queueItemURL != fmt.Sprintf("/queue/item/%v", tc.item) {
				return 0, nil
			}
			return tc.buildID, nil
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

		for _, expectedMsg := range tc.messages {
			t.Log(expectedMsg.Message)
			msg := <-channel
			if msg.Message != expectedMsg.Message {
				t.Fatalf("Unexpected Message in reply: %s != %s", msg.Message, expectedMsg.Message)
			}
			if msg.Error != expectedMsg.Error {
				t.Fatalf("Unexpected Error in reply: %v != %v", msg.Error, expectedMsg.Error)
			}
			if msg.Done != expectedMsg.Done {
				t.Fatalf("Unexpected Done in reply: %v != %v", msg.Done, expectedMsg.Done)
			}
		}
	}
}
