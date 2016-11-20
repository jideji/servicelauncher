package main

import (
	"fmt"
	"github.com/jideji/servicelauncher/config"
	"github.com/jideji/servicelauncher/service"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
	}

	action := os.Args[1]

	services := config.LoadServices()

	selected := resolveServices(services, os.Args[2:]...)

	for _, service := range selected {
		doAction(service, action)
	}
}

func resolveServices(services map[string]*service.Service, names ...string) []*service.Service {
	var selected []*service.Service
	if len(names) > 0 {
		for _, name := range names {
			service, ok := services[name]
			if !ok {
				println(fmt.Sprintf("No service named '%s' found.", name))
				os.Exit(10)
			}
			selected = append(selected, service)
		}
	} else {
		for _, service := range services {
			selected = append(selected, service)
		}
	}

	return selected
}

func doAction(service *service.Service, action string) {
	running := service.IsRunning()

	if action == "status" {
		if running {
			pid, _ := service.Pid()
			fmt.Printf("Service '%s' is running with pid %d.\n", service.Name, pid)
		} else {
			fmt.Printf("Service '%s' is not running.\n", service.Name)
		}
		return
	}

	if action == "stop" || action == "restart" {
		if running {
			pid, err := service.Pid()
			if err != nil {
				panic(err)
			}
			fmt.Printf("Killing process %d.\n", pid)
			service.Stop()
			running = false
		} else {
			fmt.Println("Not running.")
		}
	}

	if action == "start" || action == "restart" {
		if running {
			println(fmt.Sprintf("Service '%s' already running. Try restart.", service.Name))
			os.Exit(11)
		}
		p, err := service.Start()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Service '%s' started with pid %d.\n", service.Name, p.Pid)
	}
}

func showHelp() {
	fmt.Fprint(os.Stderr, "SYNTAX:\n")
	fmt.Fprintf(os.Stderr, "\t%s <action> [<service name>]\n", os.Args[0])
	fmt.Fprint(os.Stderr, "Actions:\n")
	fmt.Fprint(os.Stderr, "\tstart, stop, restart, status\n")
	fmt.Fprint(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "\t%s start httpserver\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s status\n", os.Args[0])
	os.Exit(1)
}
