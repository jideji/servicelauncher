package config

import (
	"github.com/jideji/servicelauncher/service"
	"github.com/magiconair/properties"
	"os"
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(`^(service\.([^.]+))\.command`)

// LoadServices loads service definitions from the config file ($HOME/.slcfg).
// It will panic if the file does not exist.
func LoadServices() *service.Services {
	p := properties.MustLoadFile(os.ExpandEnv("$HOME/.slcfg"), properties.UTF8)

	var services []service.Service
	for _, key := range p.Keys() {
		if regex.MatchString(key) {
			submatch := regex.FindStringSubmatch(key)
			prefix := submatch[1]
			name := submatch[2]
			command := p.MustGetString(prefix + ".command")
			commandPattern := p.MustGetString(prefix + ".pattern")
			directory := p.GetString(prefix+".directory", "")
			labelsStr := p.GetString(prefix+".labels", "")
			labels := parseLabels(labelsStr)

			srv := service.NewExternalService(
				name,
				commandPattern,
				command,
				labels,
				directory)
			services = append(services, srv)
		}
	}

	return service.NewServices(services)
}

func parseLabels(labelsStr string) []string {
	var labels []string
	splitLabels := strings.Split(labelsStr, ",")
	for _, label := range splitLabels {
		trimmed := strings.TrimSpace(label)
		if len(trimmed) > 0 {
			labels = append(labels, trimmed)
		}
	}
	return labels
}
