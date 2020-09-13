package jenkins

import (
	"time"

	"github.com/bndr/gojenkins"
)

// Build ...
type Build struct {
	*gojenkins.Build
}

// NewBuild ...
func NewBuild(build *gojenkins.Build) *Build {
	return &Build{
		build,
	}
}

// Poll ...
func (b *Build) Poll() (int, error) {
	if b.Build.Jenkins.Version != "" {
		return b.Build.Poll()
	}
	b.Raw.Building = false
	return 200, nil
}

// Wait ...
func (b *Build) Wait() (err error) {
	for {
		if !b.Raw.Building {
			break
		}
		time.Sleep(100 * time.Millisecond)
		_, err = b.Poll()
		if err != nil {
			return err
		}
	}
	return
}
