package main

import (
	"bytes"
	"encoding/json"
	"github.com/jideji/servicelauncher/service"
	"github.com/jideji/servicelauncher/web"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestReturnsEmptyListWhenNoServicesConfigured(t *testing.T) {
	srv := httptest.NewServer(web.WebHandler(service.Services{}))
	defer srv.Close()

	actual := getList(t, srv.URL+"/api")

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

	actual := getList(t, srv.URL+"/api")

	expected := []web.ServiceStatus{
		{Name: "name1", Status: "running"},
		{Name: "name2", Status: "stopped"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Got %s ; Wanted %s", expected, actual)
	}
}

func TestStartsService(t *testing.T) {
	fake := FakeService{"name", false}
	srv := httptest.NewServer(web.WebHandler(service.Services{
		"name": &fake,
	}))
	defer srv.Close()

	post(t, srv.URL+"/api/name/start")

	if !fake.running {
		t.Error("Expected service to be running")
	}
}

func getList(t *testing.T, url string) []web.ServiceStatus {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %s", err.Error())
	}

	d := json.NewDecoder(bytes.NewReader(b))

	gs := &web.GetServices{}
	err = d.Decode(gs)
	if err != nil {
		t.Fatalf("Failed decoding response: %s\nResponse:%s", err.Error(), string(b))
	}

	return gs.Services
}

func post(t *testing.T, url string) {
	resp, err := http.Post(url, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		t.Fatalf("Got status code %d; Wanted 202", resp.StatusCode)
	}
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
	s.running = true
	return nil
}
func (s *FakeService) Stop() error {
	s.running = false
	return nil
}
