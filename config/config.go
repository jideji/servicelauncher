package config

import (
	"github.com/jideji/servicelauncher/service"
	"github.com/magiconair/properties"
	"os"
	"regexp"
)

var regex = regexp.MustCompile(`^(service\.([^.]+))\.command`)

func LoadServices() map[string]*service.Service {
	p := properties.MustLoadFile(os.ExpandEnv("$HOME/.slcfg"), properties.UTF8)

	services := make(map[string]*service.Service)
	for _, key := range p.Keys() {
		if regex.MatchString(key) {
			submatch := regex.FindStringSubmatch(key)
			prefix := submatch[1]
			name := submatch[2]
			command := p.MustGetString(prefix + ".command")
			commandPattern := p.MustGetString(prefix + ".pattern")

			srv := service.Service{
				Name:    name,
				Pattern: commandPattern,
				Command: command,
			}
			services[name] = &srv
		}
	}

	return services
}
