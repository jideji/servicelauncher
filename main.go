package main

import (
	"fmt"
	"github.com/jideji/servicelauncher/config"
	"github.com/jideji/servicelauncher/service"
	"github.com/jideji/servicelauncher/web"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
	}

	action := os.Args[1]

	services := config.LoadServices()

	if action == "server" {
		h := web.WebHandler(services)
		http.Handle("/", h)
		fmt.Println("Listening on port :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
		return
	}

	selected := resolveServices(services, os.Args[2:]...)

	lastErrCode := 0
	for _, service := range selected {
		err := doAction(service, action)
		if err != nil {
			if c, ok := err.(CmdError); ok {
				println(c.Msg)
				lastErrCode = c.Code
			} else {
				println(err.Error())
			}
		}
	}

	os.Exit(lastErrCode)
}

func resolveServices(services map[string]service.Service, names ...string) []service.Service {
	var selected []service.Service
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

func doAction(srv service.Service, action string) error {
	running, err := srv.IsRunning()
	if err != nil {
		return err
	}

	if action == "status" {
		if running {
			pid, _ := srv.Pid()
			fmt.Printf("Service '%s' is running with pid %d.\n", srv.Name(), pid)
		} else {
			fmt.Printf("Service '%s' is not running.\n", srv.Name())
		}
		return nil
	}

	if action == "stop" || action == "restart" {
		if running {
			pid, err := srv.Pid()
			if err != nil {
				return err
			}
			fmt.Printf("Killing service '%s' (process %d).\n", srv.Name(), pid)
			srv.Stop()
			running = false
		} else {
			fmt.Printf("Service '%s' not running.\n", srv.Name())
		}
	}

	if action == "start" || action == "restart" {
		if running {
			return cmdError(fmt.Sprintf("Service '%s' already running. Try restart.", srv.Name()), 11)
		}
		err := srv.Start()
		if err != nil {
			return err
		}
		pid, err := srv.Pid()
		if err != nil {
			return err
		}
		fmt.Printf("Service '%s' started with pid %d.\n", srv.Name(), pid)
	}
	return nil
}

func showHelp() {
	fmt.Fprint(os.Stderr, "SYNTAX:\n")
	fmt.Fprintf(os.Stderr, "\t%s <action> [<service name>]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s server\n", os.Args[0])
	fmt.Fprint(os.Stderr, "Actions:\n")
	fmt.Fprint(os.Stderr, "\tstart, stop, restart, status\n")
	fmt.Fprint(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "\t%s start httpserver\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s status\n", os.Args[0])
	os.Exit(1)
}

func cmdError(msg string, code int) error {
	return CmdError{Msg: msg, Code: code}
}

type CmdError struct {
	Msg  string
	Code int
}

func (c CmdError) Error() string {
	return c.Msg
}
