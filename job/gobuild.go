package job

import (
	"errors"
	"fmt"

	"github.com/barryz/goci/config"
	"github.com/barryz/goci/util"
)

// GoBuildJob represents a job with go build.
type GoBuildJob string

// Do do an actual job.
func (j GoBuildJob) Do() (msg string, err error) {
	buildCmd := config.DefaultConfig.Build

	fmt.Println(">>", buildCmd)
	out, err := util.ShExec(buildCmd, 400)
	fmt.Print(out)
	if err != nil {
		return "", errors.New("go build error" + err.Error())
	}
	return "", nil
}

// IsFailTerminate indicates whether to terminate when job execute fails.
func (j GoBuildJob) IsFailTerminate() bool {
	return true
}

// Name returns the job name.
func (j GoBuildJob) Name() string {
	return string(j)
}
