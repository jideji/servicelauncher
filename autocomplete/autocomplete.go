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
local -a options args current_arg
# All arguments excluding the command
args=($words)
args[1]=()
# Where are we (excluding command)
current_arg=$[ $CURRENT - 1 ]
# call servicelauncher to resolve auto-complete candidates
options=("${(@0)$(servicelauncher "$current_arg" "$PREFIX" $args --autocomplete-options)}")
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
	// Remove autocomplete flag and anything coming after
	ac := indexOfAutocomplete(args...)
	if ac != -1 {
		args = args[:ac]
	}

	// Remove prefix if present
	if len(prefix) > 0 {
		args = args[0 : position-1]
	}

	// List configured services
	// (No matching required - the shell will do that for us)
	if position > 1 {
		args = args[1:]
		services := serviceLoader()

		nameSet := make(set)
		for _, srv := range services.AsSlice() {
			if !contains(args, srv.Name()) {
				nameSet.Add(srv.Name())
			}
		}
		names := nameSet.AsSortedSlice()

		labelSet := make(set)
		for _, srv := range services.AsSlice() {
			for _, label := range srv.Labels() {
				if !contains(args, fmt.Sprintf("l:%s", label)) {
					labelSet.Add(fmt.Sprintf("l\\:%s", label))
				}
			}
		}
		labels := labelSet.AsSortedSlice()
		for _, label := range labels {
			names = append(names, label)
		}

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
