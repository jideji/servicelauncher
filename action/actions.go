package action

import (
	"fmt"
	"github.com/jideji/servicelauncher/service"
	"os"
)

// StatusCode a program exit code
type StatusCode int

const (
	success = StatusCode(iota)
	warning
	failure
	seriousFailure
)

// escalate returns whichever status code is the most serious
func (s StatusCode) escalate(other StatusCode) StatusCode {
	if other > s {
		return other
	}
	return s
}

// Action represents a command that can be performed on services
type Action interface {
	Name() string
	Description() string
	Perform(services ...service.Service) StatusCode
}

var actions = map[string]Action{
	"list":    &listAction{},
	"status":  &statusAction{},
	"start":   &startAction{},
	"stop":    &stopAction{},
	"restart": &restartAction{},
}

// FindAction finds a command by name
func FindAction(name string) Action {
	return actions[name]
}

type listAction struct{}

func (la *listAction) Name() string        { return "list" }
func (la *listAction) Description() string { return "List available services" }
func (la *listAction) Perform(services ...service.Service) StatusCode {
	for _, srv := range services {
		fmt.Println(srv.Name())
	}
	return success
}

type statusAction struct{}

func (la *statusAction) Name() string        { return "status" }
func (la *statusAction) Description() string { return "Check status of services" }
func (la *statusAction) Perform(services ...service.Service) StatusCode {
	worstCode := success
	for _, srv := range services {
		running, err := srv.IsRunning()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", srv.Name(), err.Error())
			worstCode = worstCode.escalate(seriousFailure)
		}

		if running {
			pid, _ := srv.Pid()
			fmt.Printf("(\033[32mx\033[0m) Service '%s' is running with pid %d.\n", srv.Name(), pid)
		} else {
			fmt.Printf("( ) Service '%s' is not running.\n", srv.Name())
		}
	}
	return worstCode
}

type startAction struct{}

func (la *startAction) Name() string        { return "start" }
func (la *startAction) Description() string { return "Start services" }
func (la *startAction) Perform(services ...service.Service) StatusCode {
	worstCode := success
	for _, srv := range services {
		statusCode := start(srv)
		worstCode = worstCode.escalate(statusCode)
	}
	return worstCode
}

func start(srv service.Service) StatusCode {
	running, err := srv.IsRunning()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %s\n", srv.Name(), err.Error()))
		return seriousFailure
	}

	if running {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("(\033[31m!\033[0m) Service '%s' already running. Try restart.", srv.Name()))
		return warning
	}

	err = srv.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %s\n", srv.Name(), err.Error()))
		return seriousFailure
	}

	pid, err := srv.Pid()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %s\n", srv.Name(), err.Error()))
		return seriousFailure
	}

	fmt.Printf("Service '%s' started with pid %d.\n", srv.Name(), pid)
	return success
}

type stopAction struct{}

func (la *stopAction) Name() string        { return "stop" }
func (la *stopAction) Description() string { return "Stop services" }
func (la *stopAction) Perform(services ...service.Service) StatusCode {
	worstCode := success
	for _, srv := range services {
		statusCode := stop(srv)
		worstCode = worstCode.escalate(statusCode)
	}
	return worstCode
}

func stop(srv service.Service) StatusCode {
	running, err := srv.IsRunning()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", srv.Name(), err.Error())
		return seriousFailure
	}

	if !running {
		fmt.Printf("(\033[33m!\033[0m) Service '%s' not running.\n", srv.Name())
		return success
	}

	pid, err := srv.Pid()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", srv.Name(), err.Error())
		return seriousFailure
	}

	fmt.Printf("Killing service '%s' (process %d).\n", srv.Name(), pid)
	err = srv.Stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", srv.Name(), err.Error())
		return seriousFailure
	}

	return success
}

type restartAction struct{}

func (la *restartAction) Name() string        { return "restart" }
func (la *restartAction) Description() string { return "Restart services" }
func (la *restartAction) Perform(services ...service.Service) StatusCode {
	worstCode := success
	for _, srv := range services {
		statusCode := stop(srv)

		if statusCode > warning {
			continue
		}

		statusCode = start(srv)
		worstCode = worstCode.escalate(failure)
	}
	return worstCode
}
