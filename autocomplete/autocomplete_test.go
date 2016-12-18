package autocomplete

import (
	"github.com/jideji/servicelauncher/service"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

const noPrefix = ""

func TestIsAutoCompleteWhenFlagSet(t *testing.T) {
	assert.Equal(t, 1, indexOfAutocomplete("arg", "--autocomplete-options"), "index of flag")
}

func TestIsNotAutoCompleteWhenMissingFlag(t *testing.T) {
	assert.Equal(t, -1, indexOfAutocomplete("arg"))
}

func TestAutoCompleteCommandsWithoutPrefix(t *testing.T) {
	results := autocomplete(
		serviceLoaderThatShouldNotBeCalled(t),
		1,
		noPrefix)

	sort.Strings(results)
	assert.Equal(t,
		[]string{
			"list:List available services",
			"restart:Restart services",
			"start:Start services",
			"status:Check status of services",
			"stop:Stop services"},
		results,
		"expected commands")
}

func TestAutoCompleteCommandsWithPrefix(t *testing.T) {
	results := autocomplete(
		serviceLoaderThatShouldNotBeCalled(t),
		1,
		"lis",
		"lis")

	sort.Strings(results)
	assert.Equal(t,
		[]string{
			"list:List available services",
			"restart:Restart services",
			"start:Start services",
			"status:Check status of services",
			"stop:Stop services"},
		results,
		"expected commands")
}

func TestAutoCompleteServices(t *testing.T) {
	results := autocomplete(
		serviceLoader(srv("webserver"), srv("http-proxy")),
		2,
		noPrefix,
		"status")

	assert.Equal(t,
		[]string{"http-proxy", "webserver"},
		results,
		"expected commands")
}

func TestAutoCompleteServiceLabels(t *testing.T) {
	results := autocomplete(
		serviceLoader(srv("webserver", "group2"), srv("http-proxy", "group1", "group2")),
		2,
		noPrefix, "")

	assert.Equal(t,
		[]string{"http-proxy", "webserver", "l\\:group1", "l\\:group2"},
		results,
		"expected commands")
}

func TestSkipsServicesAlreadyGiven(t *testing.T) {
	results := autocomplete(
		serviceLoader(srv("webserver"), srv("http-proxy")),
		2,
		noPrefix,
		"status", "webserver")

	assert.Equal(t,
		[]string{"http-proxy"},
		results,
		"expected commands")
}

func TestIgnoresEntriesAfterAutocompleteFlag(t *testing.T) {
	results := autocomplete(
		serviceLoaderThatShouldNotBeCalled(t),
		1,
		noPrefix,
		"--autocomplete-options", "status")

	sort.Strings(results)
	assert.Equal(t,
		[]string{
			"list:List available services",
			"restart:Restart services",
			"start:Start services",
			"status:Check status of services",
			"stop:Stop services"},
		results,
		"expected commands")
}

func serviceLoaderThatShouldNotBeCalled(t *testing.T) service.Loader {
	return func() *service.Services {
		t.Error("Service loader called")
		return service.NewServices([]service.Service{})
	}
}

func serviceLoader(services ...service.Service) service.Loader {
	var s []service.Service
	for _, srv := range services {
		s = append(s, srv)
	}
	return func() *service.Services {
		return service.NewServices(s)
	}
}

func srv(name string, labels ...string) service.Service {
	return &DummyService{name, labels}
}

type DummyService struct {
	name   string
	labels []string
}

func (s *DummyService) Start() error             { return nil }
func (s *DummyService) Pid() (int, error)        { return -1, nil }
func (s *DummyService) Name() string             { return s.name }
func (s *DummyService) Labels() []string         { return s.labels }
func (s *DummyService) IsRunning() (bool, error) { return false, nil }
func (s *DummyService) Stop() error              { return nil }
