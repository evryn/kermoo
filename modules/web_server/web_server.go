package web_server

import (
	"context"
	"fmt"
	"kermoo/config"
	"kermoo/modules/fluent"
	"kermoo/modules/logger"
	"kermoo/modules/planner"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"go.uber.org/zap"
)

type WebServerFault struct {
	PlanRefs []string `json:"planRefs"`

	Percentage fluent.FluentFloat     `json:"percentage"`
	Interval   *fluent.FluentDuration `json:"interval"`
	Duration   *fluent.FluentDuration `json:"duration"`
}

type WebServer struct {
	planner.CanAssignPlan
	Routes      []*Route        `json:"routes"`
	Interface   *string         `json:"interface"`
	Port        *int32          `json:"port"`
	Fault       *WebServerFault `json:"fault"`
	server      *http.Server
	isListening bool
}

func (ws *WebServer) GetName() string {
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
		methods, _ := route.GetMethods()
		r.HandleFunc(route.Path, route.Handle).Methods(methods...)
	}

	ws.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ws.GetInterface(), ws.GetPort()),
		Handler: r,
	}

	go func() {
		logger.Log.Info("listening webserver...", zap.String("webserver", ws.GetName()))

		ws.isListening = true
		if err := ws.server.ListenAndServe(); err != nil {
			ws.isListening = false

			if err != http.ErrServerClosed {
				logger.Log.Fatal(
					"failed on listening and serving",
					zap.Error(err),
					zap.String("address", ws.server.Addr),
				)
			} else {
				logger.Log.Info("webserver is down", zap.String("webserver", ws.GetName()), zap.NamedError("reason", err))
			}
		}
	}()

	return nil
}

func (ws *WebServer) Stop() error {
	logger.Log.Info("shutting down webserver...", zap.String("webserver", ws.GetName()))
	if ws.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	return ws.server.Shutdown(ctx)
}

func (ws *WebServer) HasInlinePlan() bool {
	return ws.MakeInlinePlan() != nil
}

func (ws *WebServer) MakeInlinePlan() *planner.Plan {
	if ws.Fault == nil {
		return nil
	}

	plan := planner.NewPlan(planner.Plan{
		Percentage: &ws.Fault.Percentage,
		Interval:   ws.Fault.Interval,
		Duration:   ws.Fault.Duration,
	})

	return &plan
}

// Create a lifetime-long plan to serve webserver
func (ws *WebServer) MakeDefaultPlan() *planner.Plan {
	plan := planner.NewPlan(planner.Plan{})

	// Value of 0.0 indicates that the webserver will never fail.
	plan.Percentage = fluent.NewMustFluentFloat("0.0")

	return &plan
}

func (ws *WebServer) GetDesiredPlanNames() []string {
	if ws.Fault == nil {
		return nil
	}

	return ws.Fault.PlanRefs
}

func (ws *WebServer) getPlanPercentageState() bool {
	shouldListen := true

	for _, plan := range ws.GetAssignedPlans() {
		if !*plan.GetCurrentValue().ComputedPercentageChance {
			shouldListen = false
			break
		}
	}

	return shouldListen
}

func (ws *WebServer) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		shouldListen := ws.getPlanPercentageState()

		if shouldListen && !ws.isListening {
			if err := ws.ListenOnBackground(); err != nil {
				logger.Log.Error("error while listening to webserver", zap.Error(err))
			}
		} else if !shouldListen && ws.isListening {
			if err := ws.Stop(); err != nil {
				logger.Log.Error("error while stopping webserver", zap.Error(err))
			}
		}

		return planner.PLAN_SIGNAL_CONTINUE
	})

	return planner.CycleHooks{
		PreSleep: &preSleep,
	}
}
