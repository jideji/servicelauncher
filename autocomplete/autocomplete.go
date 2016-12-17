package autocomplete

import (
	"fmt"
	"github.com/jideji/servicelauncher/action"
	"github.com/jideji/servicelauncher/service"
	"os"
	"sort"
	"strings"
)

// IsAutocomplete returns true if --autocomplete-options can be found
// among arguments
func IsAutocomplete() bool {
	return indexOfAutocomplete(os.Args...) != -1
}

// Autocomplete returns autocompletion candidates for given level
func Autocomplete(serviceLoader service.Loader) {
	options := autocomplete(serviceLoader, os.Args[1], os.Args[2:]...)

	fmt.Println(strings.Join(options, "\x00"))

	os.Exit(0)
}

func indexOfAutocomplete(args ...string) int {
	for i, arg := range args {
		if arg == "--autocomplete-options" {
			return i
		}
	}
	return -1
}

func autocomplete(serviceLoader service.Loader, prefix string, args ...string) []string {
	// Remove autocomplete flag and anything coming after
	ac := indexOfAutocomplete(args...)
	if ac != -1 {
		args = args[:ac]
	}

	// Remove prefix if present
	if len(prefix) > 0 {
		args = args[0 : len(args)-1]
	}

	// List configured services
	// (No matching required - the shell will do that for us)
	if len(args) >= 1 {
		args = args[1:]
		services := serviceLoader()
		var names []string
		for name := range services {
			if !contains(args, name) {
				names = append(names, name)
			}
		}
		sort.Strings(names)
		return names
	}

	// List commands
	var commands []string
	for _, a := range action.All() {
		commands = append(commands, fmt.Sprintf("%s:%s", a.Name(), a.Description()))
	}
	return commands
}

func contains(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}
	return false
}
