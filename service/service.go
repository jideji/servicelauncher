package service

import (
	"errors"
	"fmt"
	"github.com/jideji/servicelauncher/procs"
	"os"
	"os/exec"
	"syscall"
)

// Service represents a service that can be started.
type Service struct {
	Name    string
	Pattern string
	Command string
	process *procs.Process
}

// Start runs the service using the service command.
// It redirects stdout+stderr to /tmp/<servicename>.log.
func (s *Service) Start() (*procs.Process, error) {
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

// Pid returns the process id of the running service.
// Returns an error if process is not running.
func (s *Service) Pid() (int, error) {
	p := s.getProcess()
	if p == nil {
		return -1, errors.New("No process found.")
	}
	return p.Pid, nil
}

// IsRunning returns true if process is running.
func (s *Service) IsRunning() bool {
	return s.getProcess() != nil
}

// Stop kills the running process.
func (s *Service) Stop() {
	p := s.getProcess()
	if err := p.Kill(); err != nil {
		panic(err)
	}
	s.process = nil
}

func (s *Service) getProcess() *procs.Process {
	if s.process == nil {
		pr, err := procs.FindByCommandLine(s.Pattern)
		if err != nil {
			panic(err)
		}
		s.process = pr
	}
	return s.process
}
