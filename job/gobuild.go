package job

import (
	"errors"
	"fmt"

	"github.com/barryz/goci/config"
	"github.com/barryz/goci/util"
)

type GoBuildJob string

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

func (j GoBuildJob) IsFailTerminate() bool {
	return true
}

func (j GoBuildJob) Name() string {
	return string(j)
}
