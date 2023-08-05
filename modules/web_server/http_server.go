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
	"github.com/gosimple/slug"
	"go.uber.org/zap"
)

type WebServer struct {
	planner.PlannableTrait
	Routes        []*Route       `json:"routes"`
	Interface     *string        `json:"interface"`
	Port          *int32         `json:"port"`
	InitiateAfter *time.Duration `json:"initiate_after"`
	Plan          *planner.Plan  `json:"plan"`
	PlanRefs      []string       `json:"plan_refs"`
	server        *http.Server
}

func (ws *WebServer) GetUid() string {
	return slug.Make(fmt.Sprintf("webserver-%s-%d", ws.GetInterface(), ws.GetPort()))
}

func (ws *WebServer) GetPort() int32 {
	if ws.Port != nil {
		return *ws.Port
	}

	return config.Default.WebServer.Port
}

func (ws *WebServer) GetInterface() string {
	if ws.Interface != nil {
		return *ws.Interface
	}

	return config.Default.WebServer.Interface
}

func (ws *WebServer) Validate() error {
	if ws.Routes == nil {
		return fmt.Errorf("no routes are provided")
	}

	return nil
}

func (ws *WebServer) ListenOnBackground() error {
	invalid := ws.Validate()

	if invalid != nil {
		return invalid
	}

	r := mux.NewRouter()

	for _, route := range ws.Routes {
		r.HandleFunc(route.Path, route.Handle).Methods(route.Methods...)
	}

	ws.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ws.GetInterface(), ws.GetPort()),
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

func (ws *WebServer) HasCustomPlan() bool {
	return ws.Plan != nil
}

func (ws *WebServer) MakeCustomPlan() *planner.Plan {
	plan := *ws.Plan
	return &plan
}

func (ws *WebServer) GetDesiredPlanNames() []string {
	return ws.PlanRefs
}

func (ws *WebServer) GetPlanCallbacks() planner.Callbacks {
	return planner.Callbacks{
		PreSleep: func(ep *planner.ExecutablePlan, ev *planner.ExecutableValue) planner.PlanSignal {
			return planner.PLAN_SIGNAL_CONTINUE
		},
		PostSleep: func(startedAt time.Time, timeSpent time.Duration) planner.PlanSignal {
			return planner.PLAN_SIGNAL_TERMINATE
		},
	}
}
