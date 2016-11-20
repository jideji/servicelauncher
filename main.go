package main

import (
	"github.com/jideji/servicelauncher/procs"
	"github.com/jideji/servicelauncher/props"
)

func main() {
	services := props.LoadServices()

	for name, service := range services {
		pr, err := procs.FindByCommandLine(service.Command)
		if err != nil {
			panic(err)
		}

		println(name + ":")
		println("  " + pr.Pid)
		println("  " + pr.CommandLine)
	}
}
