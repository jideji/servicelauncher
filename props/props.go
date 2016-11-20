package props

import (
	"fmt"
	"github.com/jideji/servicelauncher/procs"
	"github.com/magiconair/properties"
	"os"
	"os/exec"
	"regexp"
	"syscall"
)

type Service struct {
	Name    string
	Pattern string
	Command string
}

var regex = regexp.MustCompile(`^(service\.([^.]+))\.command`)

func LoadServices() map[string]*Service {
	p := properties.MustLoadFile(os.ExpandEnv("$HOME/.slcfg"), properties.UTF8)

	services := make(map[string]*Service)
	for _, key := range p.Keys() {
		if regex.MatchString(key) {
			submatch := regex.FindStringSubmatch(key)
			prefix := submatch[1]
			name := submatch[2]
			command := p.MustGetString(prefix + ".command")
			commandPattern := p.MustGetString(prefix + ".pattern")

			srv := Service{
				Name:    name,
				Pattern: commandPattern,
				Command: command,
			}
			services[name] = &srv
		}
	}

	return services
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
