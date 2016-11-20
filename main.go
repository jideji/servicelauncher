package main

import (
	"fmt"
	"github.com/jideji/servicelauncher/procs"
	"github.com/jideji/servicelauncher/props"
	"os"
)

func main() {
	action := os.Args[1]
	serviceName := os.Args[2]

	services := props.LoadServices()

	service := services[serviceName]
	if service == nil {
		println(fmt.Sprintf("No service '%s' found", serviceName))
		os.Exit(1)
	}

	pr, err := procs.FindByCommandLine(service.Pattern)
	if err != nil {
		panic(err)
	}

	if action == "start" {
		if pr != nil {
			println(fmt.Sprintf("Service '%s' already running. Try restart.", service.Name))
			os.Exit(10)
		}
		p, err := service.Start()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Service '%s' started with pid %d.\n", service.Name, p.Pid)
	}
}
