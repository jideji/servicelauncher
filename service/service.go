package service

import (
	"fmt"
	"github.com/jideji/servicelauncher/procs"
	"os"
	"os/exec"
	"syscall"
)

type Service struct {
	Name    string
	Pattern string
	Command string
}

func (s Service) Start() (*procs.Process, error) {
	logfile, err := os.Create(fmt.Sprintf("/tmp/%s.log", s.Name))
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("bash", "-c", s.Command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = logfile
	cmd.Stderr = logfile

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	p, err := procs.FindByPid(cmd.Process.Pid)
	if err != nil {
		return nil, err
	}

	return p, nil
}
