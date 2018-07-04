package job

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/barryz/goci/config"
	"github.com/barryz/goci/util"
)

var (
	re = regexp.MustCompile("coverage: (.+?)%")
)

type GoTestJob struct {
	name                string
	packagesWithoutTest []string
	success             bool
}

func NewTestJob(name string) *GoTestJob {
	return &GoTestJob{
		name: name,
	}
}

func (j *GoTestJob) Do() (msg string, err error) {
	defer func() {
		if err == nil {
			j.success = true
		}
	}()

	testCmd := config.DefaultConfig.Test

	fmt.Println(">>", testCmd)
	out, err := util.ShExec(testCmd, 600)
	fmt.Print(out)

	lines := strings.Split(out, "\n")
	totalPckage := 0
	totalPercent := float64(0)
	for _, line := range lines {
		if strings.HasPrefix(line, "FAIL") || strings.HasPrefix(line, "--- FAIL:") {
			err = errors.New("Has FAIL")
		}
		if strings.Index(line, "coverage:") != -1 {
			totalPckage++
			matchs := re.FindStringSubmatch(line)
			if len(matchs) != 2 {
				continue
			}
			percent, _ := strconv.ParseFloat(matchs[1], 64)
			totalPercent += percent
		} else if strings.HasSuffix(line, "[no test files]") {
			j.packagesWithoutTest = append(j.packagesWithoutTest,
				strings.Trim(strings.TrimSuffix(line, "[no test files]"), "? \t"),
			)
		}
	}
	return "", err
}

func (j *GoTestJob) IsFailTerminate() bool {
	return false
}

func (j *GoTestJob) Name() string {
	return string(j.name)
}

func Test() {
	fmt.Println("test for test 01")
	fmt.Println("test for test 02")
	fmt.Println("test for test 03")
	fmt.Println("test for test 04")
	fmt.Println("test for test 05")
	fmt.Println("test for test 06")
	fmt.Println("test for test 07")
	fmt.Println("test for test 08")
	fmt.Println("test for test 09")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
	fmt.Println("test for test 02")
}
