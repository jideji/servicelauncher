package service

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Services is a map of service name to service
type Services struct {
	byName  map[string]Service
	byLabel map[string][]Service
}

func NewServices(services []Service) *Services {
	s := Services{
		byName:  make(map[string]Service),
		byLabel: make(map[string][]Service),
	}

	for _, srv := range services {
		s.byName[srv.Name()] = srv
		for _, label := range srv.Labels() {
			s.byLabel[label] = append(s.byLabel[label], srv)
		}
	}

	return &s
}

// Loader is a function for loading Services
type Loader func() *Services

// Service represents a program that can be controlled
type Service interface {
	IsRunning() (bool, error)
	Name() string
	Labels() []string
	Pid() (int, error)
	Start() error
	Stop() error
}

// AsSlice returns the given services, sorted by name
// If no service names are given, all are returned
func (s *Services) AsSlice(names ...string) []Service {
	var selected []Service
	if len(names) > 0 {
		for _, name := range names {
			if strings.HasPrefix(name, "l:") {
				name = name[2:]
				services, ok := s.byLabel[name]
				if !ok {
					println(fmt.Sprintf("No label named '%s' found.", name))
					os.Exit(10)
				}
				for _, service := range services {
					selected = append(selected, service)
				}
				continue
			}

			service, ok := s.byName[name]
			if !ok {
				println(fmt.Sprintf("No service named '%s' found.", name))
				os.Exit(10)
			}
			selected = append(selected, service)
		}
	} else {
		for _, service := range s.byName {
			selected = append(selected, service)
		}
	}
	sort.Sort(byName(selected))

	return selected
}

// byName implements sort.Interface for []ServiceStatus based on
// the Name field.
type byName []Service

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }
