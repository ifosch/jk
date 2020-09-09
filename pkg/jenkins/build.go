package jenkins

import (
	"time"

	"github.com/bndr/gojenkins"
)

var waitForBuild = func(build *gojenkins.Build) (err error) {
	for {
		if !build.Raw.Building {
			break
		}
		time.Sleep(100 * time.Millisecond)
		_, err = build.Poll()
		if err != nil {
			return
		}
	}
	return
}
