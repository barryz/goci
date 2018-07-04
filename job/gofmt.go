package job

import (
	"errors"
	"fmt"
	"strings"

	"github.com/barryz/goci/config"
	"github.com/barryz/goci/util"
)

type GoFmtJob string

func (j GoFmtJob) Do() (msg string, err error) {
	pkgs := config.DefaultConfig.RealPkgs()
	haveError := false
	for _, dir := range pkgs {
		dir = strings.TrimSuffix(dir, "/...")
		out, err := util.Execute(false, "gofmt", "-d", dir)
		if err != nil {
			haveError = true
			fmt.Println("gofmt Command Error:", err)
		}
		if len(out) != 0 {
			fmt.Print(out)
			haveError = true
		}
	}

	if haveError {
		return "", errors.New("gofmt job checked error")
	}
	return "", nil
}

func (j GoFmtJob) IsFailTerminate() bool {
	return false
}

func (j GoFmtJob) Name() string {
	return string(j)
}
