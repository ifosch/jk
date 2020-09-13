package jenkins

import (
	"fmt"
	"testing"
	"text/template"
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

	for _, tc := range tcs {
		j := &Jenkins{client: newJenkinsClientMock(tc.jobNameList, 0, 0)}

		channel := make(chan Message)
		defer close(channel)
		go j.List(channel)

		assertExpectedMessages(tc.messages, channel, t)
	}
}

func TestJenkins_Describe(t *testing.T) {
	tcs := []struct {
		jobNameList []string
		jobName     string
		tmpl        string
		messages    []Message
	}{
		{
			[]string{
				"job1",
				"job2",
			},
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

	for _, tc := range tcs {
		templ, _ := template.New("Job").Parse(tc.tmpl)
		j := &Jenkins{client: newJenkinsClientMock(tc.jobNameList, 0, 0)}

		channel := make(chan Message)
		defer close(channel)
		go j.Describe(tc.jobName, templ, channel)

		assertExpectedMessages(tc.messages, channel, t)
	}
}

func TestJenkins_Build(t *testing.T) {
	tcs := []struct {
		jobNameList []string
		jobName     string
		params      map[string]string
		item        int64
		buildID     int64
		messages    []Message
	}{
		{
			[]string{
				"job1",
				"job2",
			},
			"job1",
			map[string]string{},
			1000,
			100,
			[]Message{
				{
					Message: "Build queued for job job1",
					Error:   false,
					Done:    false,
				},
				{
					Message: "Build started http://mockedjenkins/job/job1/100",
					Error:   false,
					Done:    false,
				},
				{
					Message: "Build finished http://mockedjenkins/job/job1/100",
					Error:   false,
					Done:    true,
				},
			},
		},
		{
			[]string{
				"job1",
				"job2",
				"job_with_params",
			},
			"job_with_params",
			map[string]string{
				"ARG1": "value1",
			},
			2000,
			200,
			[]Message{
				{
					Message: "Build queued for job job_with_params",
					Error:   false,
					Done:    false,
				},
				{
					Message: "Build started http://mockedjenkins/job/job_with_params/200",
					Error:   false,
					Done:    false,
				},
				{
					Message: "Build finished http://mockedjenkins/job/job_with_params/200",
					Error:   false,
					Done:    true,
				},
			},
		},
	}

	for _, tc := range tcs {
		j := &Jenkins{client: newJenkinsClientMock(tc.jobNameList, tc.buildID, tc.item)}

		channel := make(chan Message)
		defer close(channel)
		go j.Build(tc.jobName, tc.params, channel)

		assertExpectedMessages(tc.messages, channel, t)
	}
}
