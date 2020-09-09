package jenkins

import "github.com/bndr/gojenkins"

type taskResponse struct {
	Executable struct {
		Number int64  `json:"number"`
		URL    string `json:"url"`
	} `json:"executable"`
}

//getQueueItemInfo is a helper function to get the corresponding queue
//item for a job build request, so the process can follow up with the
//build results. It is a hack to allow testing with mocks while
//gojenkins does not export taskResponse.
var getQueueItemInfo = func(j *Jenkins, queueItemURL string) (buildID int64, err error) {
	taskResponse := &taskResponse{}
	_, err = j.client.(*gojenkins.Jenkins).Requester.GetJSON(queueItemURL, taskResponse, nil)
	if err != nil {
		return 0, err
	}
	return taskResponse.Executable.Number, nil
}
