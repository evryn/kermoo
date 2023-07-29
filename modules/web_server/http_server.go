package web_server

import (
	"buggybox/config"
	"buggybox/modules/logger"
	"buggybox/modules/planner"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type WebServer struct {
	Routes        []Route        `json:"routes"`
	Interface     *string        `json:"interface"`
	Port          *int32         `json:"port"`
	InitiateAfter *time.Duration `json:"initiate_after"`
	Plan          *planner.Plan  `json:"plan"`
	PlanRef       *string        `json:"plan_ref"`
	server        *http.Server
}

func (ws *WebServer) ListenOnBackground() error {
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

	if ws.Port != nil {
		port = *ws.Port
	}

	for _, route := range ws.Routes {
		r.HandleFunc(route.Path, route.Handle).Methods(route.Methods...)
	}

	ws.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", interf, port),
		Handler: r,
	}

	go func() {
		if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal(
				"failed on listening and serving",
				zap.Error(err),
				zap.String("address", ws.server.Addr),
			)
		}
	}()

	return nil
}

func (ws *WebServer) Stop() error {
	if ws.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	return ws.server.Shutdown(ctx)
}
