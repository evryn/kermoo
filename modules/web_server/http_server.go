package web_server

import (
	"buggybox/config"
	"buggybox/modules/planner"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type WebServer struct {
	Routes        []Route        `json:"routes"`
	Interface     *string        `json:"interface"`
	Port          *int32         `json:"port"`
	InitiateAfter *time.Duration `json:"initiate_after"`
	Plan          *planner.Plan  `json:"plan"`
	PlanRef       *string        `json:"plan_ref"`
}

func (ws *WebServer) Listen() error {
	r := mux.NewRouter()

	var (
		interf = config.Default.WebServer.Interface
		port   = config.Default.WebServer.Port
	)

	if ws.Routes == nil {
		return fmt.Errorf("no routes are provided for web server")
	}

	if ws.Interface != nil {
		interf = *ws.Interface
	}

	for _, route := range ws.Routes {
		r.HandleFunc(route.Path, route.Handle).Methods(route.Methods...)
	}

	return http.ListenAndServe(fmt.Sprintf("%s:%d", interf, port), r)
}
