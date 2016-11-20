package main

import (
	"fmt"
	"github.com/jideji/servicelauncher/config"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "SYNTAX:\n")
		fmt.Fprintf(os.Stderr, "\t%s <action> <service name>\n", os.Args[0])
		fmt.Fprint(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "\t%s start httpserver\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t%s stop httpserver\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t%s restart httpserver\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t%s status httpserver\n", os.Args[0])
		os.Exit(1)
	}

	action := os.Args[1]
	serviceName := os.Args[2]

	services := config.LoadServices()

	service := services[serviceName]
	if service == nil {
		println(fmt.Sprintf("No service named '%s' found.", serviceName))
		os.Exit(10)
	}

	running := service.IsRunning()

	if action == "status" {
		if running {
			pid, _ := service.Pid()
			fmt.Printf("Service '%s' is running with pid %d.\n", service.Name, pid)
		} else {
			fmt.Printf("Service '%s' is not running.\n", service.Name)
		}
		os.Exit(0)
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
