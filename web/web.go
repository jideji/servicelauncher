package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jideji/servicelauncher/service"
	"net"
	"net/http"
	"sort"
)

type GetServices struct {
	Services []ServiceStatus `json:"services"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

type Web struct {
	listener net.Listener
}

func WebHandler(services service.Services) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api", getList(services)).Methods("GET")
	r.HandleFunc("/api/", getList(services)).Methods("GET")
	r.HandleFunc("/api/{servicename}/start", start(services)).Methods("POST")
	return r
}

func start(services service.Services) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["servicename"]

		service, ok := services[name]
		if !ok {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("Unknown service '%s'", name)))
			return
		}
		service.Start()

		w.WriteHeader(202)
	}
}

func getList(services service.Services) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var statuses []ServiceStatus
		for _, s := range services {
			var status string
			var errMsg string
			running, err := s.IsRunning()
			if err != nil {
				errMsg = err.Error()
			}
			if running {
				status = "running"
			} else {
				status = "stopped"
			}
			statuses = append(statuses, ServiceStatus{
				Name:   s.Name(),
				Status: status,
				Error:  errMsg,
			})
		}
		sort.Sort(byName(statuses))

		w.Header().Set("content-type", "application/json")
		enc := json.NewEncoder(w)
		err := enc.Encode(GetServices{Services: statuses})

		if err != nil {
			panic(err)
		}
	}
}

func (w *Web) Close() {
	w.listener.Close()
}

// ByName implements sort.Interface for []ServiceStatus based on
// the Name field.
type byName []ServiceStatus

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }
