package job

import (
	"errors"
	"fmt"
	"strings"

	"github.com/barryz/goci/config"
	"github.com/barryz/goci/util"
)

// GoVetJob represents a job with go vet.
type GoVetJob string

// Do do an actual job.
func (j GoVetJob) Do() (msg string, err error) {
	pkgs := config.DefaultConfig.RealPkgs()
	haveError := false
	for _, dir := range pkgs {
		dir = strings.TrimSuffix(dir, "/...")

		out, err := util.Execute(false, "go", "tool", "vet", "-printfuncs", "Info,Infof,Debug,Debugf,Warn,Warnf", dir)
		fmt.Print(out)
		if err != nil {
			haveError = true
			fmt.Println("govet Command Error:", err)
		}
	}
	if haveError {
		return "", errors.New("golint have errors")
	}
	return "", nil
}

// IsFailTerminate indicates whether to terminate when job execute fails.
func (j GoVetJob) IsFailTerminate() bool {
	return false
}

// Name returns the job name.
func (j GoVetJob) Name() string {
	return string(j)
}
