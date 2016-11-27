package main

import (
	"encoding/json"
	"github.com/jideji/servicelauncher/service"
	"github.com/jideji/servicelauncher/web"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestReturnsEmptyListWhenNoServicesConfigured(t *testing.T) {
	srv := httptest.NewServer(web.WebHandler(service.Services{}))
	defer srv.Close()

	actual := getList(t, srv.URL)

	var expected []web.ServiceStatus

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Got %s ; Wanted %s", expected, actual)
	}
}

func TestReturnsStateOfConfiguredServices(t *testing.T) {
	srv := httptest.NewServer(web.WebHandler(service.Services{
		"name1": &FakeService{"name1", true},
		"name2": &FakeService{"name2", false},
	}))
	defer srv.Close()

	actual := getList(t, srv.URL)

	expected := []web.ServiceStatus{
		{Name: "name1", Status: "running"},
		{Name: "name2", Status: "stopped"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Got %s ; Wanted %s", expected, actual)
	}
}

func getList(t *testing.T, url string) []web.ServiceStatus {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	gs := &web.GetServices{}
	err = d.Decode(gs)
	if err != nil {
		t.Fatalf("Failed decoding response: %s", err)
	}

	return gs.Services
}

type FakeService struct {
	name    string
	running bool
}

func (s *FakeService) Name() string {
	return s.name
}
func (s *FakeService) IsRunning() bool {
	return s.running
}
func (s *FakeService) Pid() (int, error) {
	return -1, nil
}
func (s *FakeService) Start() error {
	return nil
}
func (s *FakeService) Stop() error {
	return nil
}
