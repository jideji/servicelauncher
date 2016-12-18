package service

import (
	"errors"
	"fmt"
	"github.com/jideji/servicelauncher/procs"
	"os"
	"os/exec"
	"syscall"
)

// ExternalService lets services run separately.
// A regex pattern is used to find a running instance of a service.
type ExternalService struct {
	name      string
	pattern   string
	command   string
	directory string
	labels    []string
	process   *procs.Process
}

// NewExternalService creates a service.
func NewExternalService(
	name string,
	pattern string,
	cmd string,
	labels []string,
	dir string) Service {

	return &ExternalService{
		name:      name,
		pattern:   pattern,
		command:   cmd,
		labels:    labels,
		directory: dir,
	}
}

// Start runs the service using the service command.
// It redirects stdout+stderr to /tmp/<servicename>.log.
func (s *ExternalService) Start() error {
	logfile, err := os.Create(fmt.Sprintf("/tmp/%s.log", s.name))
	if err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", s.command)
	cmd.Dir = s.directory
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = logfile
	cmd.Stderr = logfile

	err = cmd.Start()
	if err != nil {
		return err
	}

	p, err := procs.FindByPid(cmd.Process.Pid)
	if err != nil {
		return err
	}
	s.process = p

	return nil
}

// Pid returns the process id of the running service.
// Returns an error if process is not running.
func (s *ExternalService) Pid() (int, error) {
	p, err := s.getProcess()
	if err != nil {
		return -1, err
	}
	if p == nil {
		return -1, errors.New("No process found.")
	}
	return p.Pid, nil
}

// Name returns the name of the service.
func (s *ExternalService) Name() string {
	return s.name
}

// Labels returns the label given to service, if any.
func (s *ExternalService) Labels() []string {
	return s.labels
}

// IsRunning returns true if process is running.
func (s *ExternalService) IsRunning() (bool, error) {
	process, err := s.getProcess()
	if err != nil {
		return false, err
	}
	return process != nil, nil
}

// Stop kills the running process.
func (s *ExternalService) Stop() error {
	p, err := s.getProcess()
	if err != nil {
		return err
	}

	if err = p.Kill(); err != nil {
		return err
	}
	s.process = nil
	return nil
}

func (s *ExternalService) getProcess() (*procs.Process, error) {
	if s.process == nil {
		pr, err := procs.FindByCommandLine(s.pattern)
		if err != nil {
			return nil, err
		}
		s.process = pr
	}
	return s.process, nil
}
