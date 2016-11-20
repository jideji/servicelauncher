package procs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Process represents a running process.
type Process struct {
	Pid         int
	CommandLine string
}

// FindByCommandLine finds a process using a regex.
// Returns an error if there is not exactly one match.
func FindByCommandLine(regex string) (*Process, error) {
	r := regexp.MustCompile(regex)

	procs, err := allProcesses()
	if err != nil {
		return nil, err
	}

	var filtered []*Process

	for _, p := range procs {
		if r.MatchString(p.CommandLine) {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) > 1 {
		return nil, fmt.Errorf("Found more than one match - found %d", len(filtered))
	}

	if len(filtered) == 0 {
		return nil, nil
	}

	return filtered[0], nil
}

// FindByPid finds a process by its pid.
func FindByPid(pid int) (*Process, error) {
	procs, err := allProcesses()
	if err != nil {
		return nil, err
	}

	for _, p := range procs {
		if p.Pid == pid {
			return p, nil
		}
	}

	return nil, nil
}

// Kill sends a kill signal to a process.
// It then waits (up to 5 seconds) for it do die.
// If it doesn't die within the time limit, an error is returned.
func (p *Process) Kill() error {
	syscall.Kill(p.Pid, syscall.SIGKILL)

	limit := time.After(5 * time.Second)
	for {
		p, err := FindByPid(p.Pid)
		if err != nil {
			return err
		}
		if p == nil {
			return nil
		}

		select {
		case <-limit:
			return fmt.Errorf("Timed out waiting for %d to die.", p.Pid)
		case <-time.Tick(100 * time.Millisecond):
		}
	}
}

func allProcesses() ([]*Process, error) {

	cmd := exec.Command("ps", "-ax", "-o", "pid", "-o", "command")
	buf := bytes.NewBuffer(nil)

	cmd.Stdout = buf
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")[1:]

	var ps []*Process
	for _, line := range lines {
		fields := strings.SplitN(strings.TrimSpace(line), " ", 2)
		if len(fields) != 2 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}

		commandline := strings.TrimSpace(fields[1])

		p := Process{pid, commandline}
		ps = append(ps, &p)
	}

	return ps, nil
}
