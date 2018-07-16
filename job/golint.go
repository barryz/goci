package job

import (
	"errors"
	"fmt"
	"strings"

	"github.com/barryz/goci/config"
	"github.com/barryz/goci/util"
)

type GoLintJob string

func (j GoLintJob) Do() (msg string, err error) {
	// fmt.Print("!!!! golint new config field \"lint -> ignore_no_comment_error\": https://github.com/barryz/goci#lint-结构字段\n\n")

	pkgs := config.DefaultConfig.RealPkgs()
	noComments := config.DefaultConfig.Lint.IgnoreNoCommentError
	haveError := false
	for _, dir := range pkgs {
		cmd := "golint " + dir
		if noComments {
			cmd = "golint " + dir + " | grep -v 'should have comment'"
		}
		out, _ := util.ShExec(cmd, 120)
		if len(out) != 0 {
			fmt.Println(">>", cmd)
			realOut := ""
			outs := strings.Split(out, "\n")
			for _, line := range outs {
				if len(line) == 0 {
					continue
				}
				if !strings.HasPrefix(line, "warning:") {
					realOut += line + "\n"
				} else {
					fmt.Println(line)
				}
			}
			if len(realOut) > 0 {
				fmt.Print(realOut)
				haveError = true
			}
		}
	}
	if haveError {
		return "", errors.New("golint have errors")
	}
	return "", nil
}

func (j GoLintJob) IsFailTerminate() bool {
	return false
}

func (j GoLintJob) Name() string {
	return string(j)
}
