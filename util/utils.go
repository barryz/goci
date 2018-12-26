package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// Carve carving individual field lines.
func Carve(line string) {
	fmt.Println("--------------------------", line, "--------------------------")
}

// CommandResp represents a response for the specified command.
type CommandResp struct {
	Output string
	Err    error
}

// StartCommand starting a command.
func StartCommand(cmds string, envs []string) (<-chan *CommandResp, *exec.Cmd, error) {
	dirName, _ := os.Getwd()
	fmt.Println("Current Dir:", dirName)

	fmt.Printf(">> %s\n", cmds)

	p := exec.Command("sh", "-c", cmds)
	p.Dir = dirName
	if envs != nil {
		p.Env = envs
	}
	p.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	ch := make(chan *CommandResp, 1)
	go func() {
		out, err := p.CombinedOutput()
		fmt.Printf("%s Output: %s\n", cmds, out)
		ch <- &CommandResp{
			Output: string(out),
			Err:    err,
		}
	}()
	return ch, p, nil
}

// Execute do an actual execution.
func Execute(show bool, cmd string, args ...string) (string, error) {
	p := exec.Command(cmd, args...)
	out, err := p.CombinedOutput()
	if err != nil || show {
		fmt.Print(">> ", cmd, " ")
		for _, v := range args {
			fmt.Printf("%v ", v)
		}
		fmt.Println("")
	}

	return string(out), err
}

type toBytes interface {
	Bytes() []byte
}

type result struct {
	out []byte
	err error
}

// ShExec do an actual shell execution.
func ShExec(cmd string, seconds int) (string, error) {
	p := exec.Command("sh", "-c", cmd)

	outCh := make(chan result, 1)
	go func() {
		out, err := p.CombinedOutput()
		outCh <- result{
			out: out,
			err: err,
		}
	}()
	select {
	case out := <-outCh:
		return string(out.out), out.err
	case <-time.After(time.Second * time.Duration(seconds)):
		err := p.Process.Kill()
		if err != nil {
			return "", err
		}
		select {
		case out := <-outCh:
			return string(out.out), errors.New("job timeout")
		case <-time.After(time.Second * 5):
			return "", errors.New("job timeout")
		}
	}
}

// StrToPerStr to per string.
func StrToPerStr(str string, usePlus bool) string {
	n, _ := strconv.Atoi(str)
	return IntToPerStr(n, usePlus)
}

// IntToPerStr to per string.
func IntToPerStr(n int, usePlus bool) string {
	flag := 1 // plus
	if n < 0 {
		n = -n
		flag = -1
	}

	a := n / 100
	b := n % 100
	str := strconv.Itoa(a) + "." + fmt.Sprintf("%02d", b) + "%"
	if flag == -1 {
		return "-" + str
	}
	if usePlus {
		return "+" + str
	}
	return str
}
