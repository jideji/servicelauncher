package props

import (
	"github.com/magiconair/properties"
	"os"
	"regexp"
	"strconv"
)

type Service struct {
	Name    string
	Pid     int
	Command string
}

var regex = regexp.MustCompile(`^(service\.([^.]+))\.pid`)

func LoadServices() map[string]*Service {
	p := properties.MustLoadFile(os.ExpandEnv("$HOME/.slcfg"), properties.UTF8)

	services := make(map[string]*Service)
	for _, key := range p.Keys() {
		if regex.MatchString(key) {
			submatch := regex.FindStringSubmatch(key)
			prefix := submatch[1]
			name := submatch[2]
			pid, err := strconv.Atoi(p.MustGetString(key))
			if err != nil {
				println("Property", key, "expected to be int")
				continue
			}
			commandline := p.MustGetString(prefix + ".command")

			srv := Service{
				Name:    name,
				Pid:     pid,
				Command: commandline,
			}
			services[name] = &srv
		}
	}

	return services
}
