package main

import (
	"fmt"
	"github.com/jideji/servicelauncher/action"
	"github.com/jideji/servicelauncher/autocomplete"
	"github.com/jideji/servicelauncher/config"
	"github.com/jideji/servicelauncher/service"
	"github.com/jideji/servicelauncher/web"
	"net/http"
	"os"
	"strings"
)

func main() {
	if autocomplete.IsAutocomplete() {
		autocomplete.Autocomplete(func() service.Services { return config.LoadServices() })
	}

	if len(os.Args) < 2 || os.Args[1] == "--help" {
		showHelp()
	}

	if os.Args[1] == "autocomplete-zsh" {
		showZshAutocompleteScript()
	}

	actionStr := os.Args[1]

	services := config.LoadServices()

	if actionStr == "server" {
		h := web.WebHandler(services)
		http.Handle("/", h)
		fmt.Println("Listening on port :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
		return
	}

	selected := services.AsSlice(os.Args[2:]...)

	a := action.FindAction(actionStr)
	if a == nil {
		fmt.Fprintf(os.Stderr, "Unknown command '%s'\n", actionStr)
		os.Exit(100)
	}

	statusCode := a.Perform(selected...)
	os.Exit(int(statusCode))
}

func showHelp() {
	fmt.Fprint(os.Stderr, "SYNTAX:\n")
	fmt.Fprintf(os.Stderr, "\t%s <action> [<service name>]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s server\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "\t\t- start a web server on port 8080")
	fmt.Fprintf(os.Stderr, "\t%s autocomplete-zsh\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "\t\t- print autocomplete-script for zsh")
	fmt.Fprintln(os.Stderr, "Actions:")
	fmt.Fprintln(os.Stderr, "\tlist, restart, start, status, stop")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintf(os.Stderr, "\t%s start httpserver\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s status\n", os.Args[0])
	os.Exit(1)
}

func showZshAutocompleteScript() {
	fmt.Println(strings.TrimSpace(`
#compdef servicelauncher
# Script to place somewhere in your fpath
# (see https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org)
local -a options args
# All arguments excluding the command
args=($words)
args[1]=()
# Call servicelauncher to resolve auto-complete candidates
options=("${(@0)$(servicelauncher "$PREFIX" $args --autocomplete-options)}")
_describe 'values' options
`))
	os.Exit(0)
}
