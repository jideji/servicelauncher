package autocomplete

import (
	"fmt"
	"github.com/jideji/servicelauncher/service"
	"os"
	"sort"
	"strings"
)

func IsAutoComplete() bool {
	return indexOfAutocomplete(os.Args...) != -1
}

func Autocomplete(serviceLoader service.ServiceLoader) {
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

func autocomplete(serviceLoader service.ServiceLoader, prefix string, args ...string) []string {
	// Remove autocomplete flag and anything coming after
	ac := indexOfAutocomplete(args...)
	if ac != -1 {
		args = args[:ac]
	}

	// Remove prefix if present
	// (The autocomplete shell filters out the matching ones)
	if len(prefix) > 0 {
		args = args[0 : len(args)-1]
	}

	// List configured services
	if len(args) >= 1 {
		services := serviceLoader()
		var names []string
		for name := range services {
			names = append(names, name)
		}
		sort.Strings(names)
		return names
	}

	// List commands
	return []string{"list", "restart", "start", "status", "stop"}
}
