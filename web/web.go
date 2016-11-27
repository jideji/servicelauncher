package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jideji/servicelauncher/service"
	"net"
	"net/http"
)

type GetServices struct {
	Services []ServiceStatus `json:"services"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Web struct {
	listener net.Listener
}

func WebHandler(services service.Services) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var statuses []ServiceStatus
		for _, s := range services {
			var status string
			if s.IsRunning() {
				status = "running"
			} else {
				status = "stopped"
			}
			statuses = append(statuses, ServiceStatus{
				Name:   s.Name(),
				Status: status,
			})
		}

		w.Header().Set("content-type", "application/json")
		enc := json.NewEncoder(w)
		err := enc.Encode(GetServices{Services: statuses})

		if err != nil {
			panic(err)
		}
	})
	return r
}

func (w *Web) Close() {
	w.listener.Close()
}
