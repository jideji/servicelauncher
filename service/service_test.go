package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanLookupByName(t *testing.T) {
	service1 := srv("the-name", "a-label")
	service2 := srv("another-name", "a-label", "different-label")

	services := NewServices([]Service{service1, service2})

	actual := services.AsSlice("the-name")

	assert.Len(t, actual, 1, "one service")
	assert.Equal(t, service1, actual[0])
}

func TestCanLookupByLabel(t *testing.T) {
	service1 := srv("name-1", "a-label")
	service2 := srv("name-2", "a-label", "different-label")
	service3 := srv("name-3", "different-label")

	services := NewServices([]Service{service1, service2, service3})

	actual := services.AsSlice("l:a-label")

	assert.Len(t, actual, 2, "two services")
	assert.Equal(t, service1, actual[0])
	assert.Equal(t, service2, actual[1])
}

func srv(name string, labels ...string) Service {
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
func (s *DummyService) String() string           { return s.name }
