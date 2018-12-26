package job

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/barryz/goci/config"
	"github.com/barryz/goci/util"
)

// GoRaceJob represents a job with go race condition testing.
type GoRaceJob string

// Do do an actual job.
func (j GoRaceJob) Do() (msg string, err error) {
	race := config.DefaultConfig.Race
	if race == nil || race.Main == "" {
		return "", errors.New("race config is empty")
	}

	var mp *exec.Cmd
	var mch <-chan *util.CommandResp
	if race.MainCMD == "" {
		mch, mp, err = util.StartCommand("go run -race "+race.Main, nil)
	} else {
		mch, mp, err = util.StartCommand(race.MainCMD, nil)
	}
	if err != nil {
		return "", errors.New("start main error: " + err.Error())
	}

	var sch <-chan *util.CommandResp
	var sp *exec.Cmd
	if race.Script != "" {
		sch, sp, err = util.StartCommand(race.Script, nil)
		if err != nil {
			return "", errors.New("start script error: " + err.Error())
		}
	}

	if race.Timeout == 0 {
		race.Timeout = 30
		fmt.Println("Use Default Timeout 30s.")
	}

	isTimeout := false
	timeout := time.After(time.Duration(race.Timeout * int(time.Second)))
	haveError := false

	select {
	case resp := <-mch:
		if resp.Err != nil {
			haveError = true
			fmt.Println("Main resp error:" + resp.Err.Error())
		}
		if strings.Index(resp.Output, "WARNING: DATA RACE") != -1 || strings.Index(resp.Output, "fatal error: concurrent map") != -1 {
			haveError = true
			fmt.Println("Have data race!\n", resp.Output)
		}
	case <-timeout:
		isTimeout = true
		fmt.Println("Main timeout!!!")
	}

	if isTimeout {
		pgid, err := syscall.Getpgid(mp.Process.Pid)
		if err == nil {
			syscall.Kill(-pgid, syscall.SIGKILL) // note the minus sign
		}

	} else if !haveError && sp != nil {
		select {
		case resp := <-sch:
			fmt.Print("Script Output:", resp.Output)
			if resp.Err != nil {
				return "", errors.New("Script resp error: " + resp.Err.Error())
			}
		case <-timeout:
			fmt.Println("Script timeout!")
			isTimeout = true
		}
	}
	if isTimeout || haveError {
		if sp != nil {
			pgid, err := syscall.Getpgid(sp.Process.Pid)
			if err == nil {
				syscall.Kill(-pgid, syscall.SIGKILL) // note the minus sign
			}
		}
		if isTimeout {
			return "", errors.New("error: main timeout")
		}
		return "", errors.New("have error")
	}
	return "", nil
}

// IsFailTerminate indicates whether to terminate when job execute fails.
func (j GoRaceJob) IsFailTerminate() bool {
	return false
}

// Name returns the job name.
func (j GoRaceJob) Name() string {
	return string(j)
}
