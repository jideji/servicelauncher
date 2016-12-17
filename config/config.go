package config

import (
	"github.com/jideji/servicelauncher/service"
	"github.com/magiconair/properties"
	"os"
	"regexp"
)

var regex = regexp.MustCompile(`^(service\.([^.]+))\.command`)

// LoadServices loads service definitions from the config file ($HOME/.slcfg).
// It will panic if the file does not exist.
func LoadServices() service.Services {
	p := properties.MustLoadFile(os.ExpandEnv("$HOME/.slcfg"), properties.UTF8)

	services := make(map[string]service.Service)
	for _, key := range p.Keys() {
		if regex.MatchString(key) {
			submatch := regex.FindStringSubmatch(key)
			prefix := submatch[1]
			name := submatch[2]
			command := p.MustGetString(prefix + ".command")
			commandPattern := p.MustGetString(prefix + ".pattern")
			directory := p.GetString(prefix+".directory", "")

			srv := service.NewExternalService(
				name,
				commandPattern,
				command,
				directory)
			services[name] = srv
		}
	}

	return services
}
