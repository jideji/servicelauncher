package autocomplete

import (
	"github.com/jideji/servicelauncher/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

const noPrefix = ""

func TestIsAutoCompleteWhenFlagSet(t *testing.T) {
	assert.Equal(t, 1, indexOfAutocomplete("arg", "--autocomplete-options"), "index of flag")
}

func TestIsNotAutoCompleteWhenMissingFlag(t *testing.T) {
	assert.Equal(t, -1, indexOfAutocomplete("arg"))
}

func TestAutoCompleteFirstLevelWithoutPrefix(t *testing.T) {
	results := autocomplete(
		serviceLoaderThatShouldNotBeCalled(t),
		noPrefix)

	assert.Equal(t,
		[]string{"list", "restart", "start", "status", "stop"},
		results,
		"expected commands")
}

func TestAutoCompleteFirstLevelWithPrefix(t *testing.T) {
	results := autocomplete(
		serviceLoaderThatShouldNotBeCalled(t),
		"lis",
		"lis")

	assert.Equal(t,
		[]string{"list", "restart", "start", "status", "stop"},
		results,
		"expected commands")
}

func TestAutoCompleteServiceLevel(t *testing.T) {
	results := autocomplete(
		serviceLoader(srv("webserver"), srv("http-proxy")),
		noPrefix,
		"status")

	assert.Equal(t,
		[]string{"http-proxy", "webserver"},
		results,
		"expected commands")
}

func TestIgnoresEntriesAfterAutocompleteFlag(t *testing.T) {
	results := autocomplete(
		serviceLoaderThatShouldNotBeCalled(t),
		noPrefix,
		"--autocomplete-options", "status")

	assert.Equal(t,
		[]string{"list", "restart", "start", "status", "stop"},
		results,
		"expected commands")
}

func serviceLoaderThatShouldNotBeCalled(t *testing.T) service.ServiceLoader {
	return func() service.Services {
		t.Error("Service loader called")
		return service.Services{}
	}
}

func serviceLoader(services ...service.Service) service.ServiceLoader {
	s := make(service.Services)
	for _, srv := range services {
		s[srv.Name()] = srv
	}
	return func() service.Services {
		return s
	}
}

func srv(name string) service.Service {
	return &DummyService{name}
}

type DummyService struct{ name string }

func (s *DummyService) Start() error             { return nil }
func (s *DummyService) Pid() (int, error)        { return -1, nil }
func (s *DummyService) Name() string             { return s.name }
func (s *DummyService) IsRunning() (bool, error) { return false, nil }
func (s *DummyService) Stop() error              { return nil }
