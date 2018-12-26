package job

import (
	"os"

	"github.com/barryz/goci/util"
)

// Job job is a interface that warps several methods for an realistic ci linter job.
type Job interface {
	Do() (msg string, err error)
	IsFailTerminate() bool
	Name() string
}

// ProjectInit a default project job for initialization.
type ProjectInit string

// Do do an actual job.
func (p ProjectInit) Do() (string, error) {
	if _, err := os.Stat("./"); os.IsNotExist(err) {
		util.Execute(true, "mkdir", "-p", "tmp")
	}
	return "", nil
}

// IsFailTerminate indicates whether to terminate when job execute fails.
func (p ProjectInit) IsFailTerminate() bool {
	return true
}

// Name returns the job name.
func (p ProjectInit) Name() string {
	return string(p)
}
