package procs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

type Process struct {
	Pid         string
	CommandLine string
}

// FindByCommandLine finds a process using a regex.
// Returns an error if there is not exactly one match.
func FindByCommandLine(regex string) (*Process, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	procs, err := allProcesses()

	var filtered []*Process

	for _, p := range procs {
		if r.MatchString(p.CommandLine) {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) != 1 {
		return nil, fmt.Errorf("Did not find exactly one match - found %d", len(procs))
	}

	return filtered[0], nil
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
		fields := strings.SplitN(line, " ", 2)
		if len(fields) != 2 {
			continue
		}

		pid := strings.TrimSpace(fields[0])
		commandline := strings.TrimSpace(fields[1])

		p := Process{pid, commandline}
		ps = append(ps, &p)
	}

	return ps, nil
}
