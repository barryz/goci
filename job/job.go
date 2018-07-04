package job

import (
	"os"

	"github.com/barryz/goci/util"
)

type Job interface {
	Do() (msg string, err error)
	IsFailTerminate() bool
	Name() string
}

type ProjectInit string

func (p ProjectInit) Do() (string, error) {
	if _, err := os.Stat("./"); os.IsNotExist(err) {
		util.Execute(true, "mkdir", "-p", "tmp")
	}
	return "", nil
}

func (p ProjectInit) IsFailTerminate() bool {
	return true
}

func (p ProjectInit) Name() string {
	return string(p)
}
