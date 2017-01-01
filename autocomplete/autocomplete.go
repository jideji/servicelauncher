package autocomplete

import (
	"fmt"
	"github.com/jideji/servicelauncher/action"
	"github.com/jideji/servicelauncher/service"
	"os"
	"sort"
	"strconv"
	"strings"
)

// ScriptFile returns a zsh file for the user to place somewhere in their fpath.
func ScriptFile() string {
	return strings.TrimSpace(`
#compdef servicelauncher
# Script to place somewhere in your fpath
# (see https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org)
local -a options
# call servicelauncher to resolve auto-complete candidates
options=("${(@0)$(servicelauncher "$CURRENT" "$PREFIX" $words --autocomplete-options)}")
_describe 'values' options
`)
}

// IsAutocomplete returns true if --autocomplete-options can be found
// among arguments
func IsAutocomplete() bool {
	return indexOfAutocomplete(os.Args...) != -1
}

// Autocomplete returns autocompletion candidates for given level
func Autocomplete(serviceLoader service.Loader) {
	position, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	options := autocomplete(serviceLoader, position, os.Args[2], os.Args[3:]...)

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

func autocomplete(serviceLoader service.Loader, position int, prefix string, args ...string) []string {
	// Remove command
	args = args[1:]
	position--

	// Remove autocomplete flag and anything coming after
	ac := indexOfAutocomplete(args...)
	if ac != -1 {
		args = args[:ac]
	}

	// Remove prefix if present
	if len(prefix) > 0 {
		remove(args, position-1)
	}

	// List configured services
	// (No matching required - the shell will do that for us)
	if position > 1 {
		args = args[1:]
		services := serviceLoader()

		candidates := byName(services, args)
		candidates = concat(candidates, byLabel(services, args))

		return candidates
	}

	// List commands
	var commands []string
	for _, a := range action.All() {
		commands = append(commands, fmt.Sprintf("%s:%s", a.Name(), a.Description()))
	}
	return commands
}

func remove(a []string, i int) {
	a[i] = a[len(a)-1]
	a = a[:len(a)-1]
}

func concat(to []string, from []string) []string {
	for _, label := range from {
		to = append(to, label)
	}
	return to
}

func byName(services *service.Services, args []string) []string {
	nameSet := make(set)
	for _, srv := range services.AsSlice() {
		if !contains(args, srv.Name()) {
			nameSet.Add(srv.Name())
		}
	}

	return nameSet.AsSortedSlice()
}

func byLabel(services *service.Services, args []string) []string {
	labelMap := make(map[string][]string)
	for _, srv := range services.AsSlice() {
		for _, label := range srv.Labels() {
			if !contains(args, fmt.Sprintf("l:%s", label)) {
				prefixedLabel := fmt.Sprintf("l\\:%s", label)
				labelMap[prefixedLabel] = append(labelMap[prefixedLabel], srv.Name())
			}
		}
	}
	var labels []string
	for labelName, services := range labelMap {
		description := strings.Join(services, ",")
		labels = append(labels, fmt.Sprintf("%s:%s", labelName, description))
	}

	sort.Strings(labels)

	return labels
}

func contains(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}
	return false
}

type set map[string]interface{}

func (s *set) AddAll(keys []string) {
	for _, key := range keys {
		s.Add(key)
	}
}

func (s *set) Add(key string) {
	(*s)[key] = nil
}

func (s *set) AsSortedSlice() []string {
	var slice []string
	for key := range *s {
		slice = append(slice, key)
	}
	sort.Strings(slice)
	return slice
}
