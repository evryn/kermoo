package web_server

import (
	"encoding/json"
	"fmt"
	"kermoo/config"
	"kermoo/modules/fluent"
	"kermoo/modules/planner"
	"kermoo/modules/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

type Route struct {
	planner.CanAssignPlan

	// Path defines the route path. Here are a few examples: "/", "/livez", "/api/v1/readyz", etc.
	Path string `json:"path"`

	// Methods is an array of method names that the route should live on them, such as "HEAD",
	// "GET", "POST", "PUT", etc.
	//
	// If no methods are defined, it'll use default ones: "HEAD", "GET", "POST"
	Methods []string `json:"methods"`

	// Content represent the response of the route. You can define either an static content or
	// enable a whoami-like response which will response few things about the application and
	// the received request.
	//
	// If nothing is set, a default "Hello World" will be responded.
	Content RouteContent `json:"content"`

	// Fault defines how the route should fail. Default is no failure.
	Fault *RouteFault `json:"fault"`
}

func (route *Route) GetName() string {
	return slug.Make(fmt.Sprintf("route-%s", route.Path))
}

func (route *Route) GetDesiredPlanNames() []string {
	if route.Fault == nil {
		return nil
	}

	return route.Fault.PlanRefs
}

func (route *Route) HasInlinePlan() bool {
	return route.MakeInlinePlan() != nil
}

func (route *Route) MakeInlinePlan() *planner.Plan {
	if route.Fault == nil {
		return nil
	}

	plan := planner.NewPlan(planner.Plan{
		Percentage: &route.Fault.Percentage,
		Interval:   route.Fault.Interval,
		Duration:   route.Fault.Duration,
	})

	return &plan
}

// Create a lifetime-long plan to serve route
func (route *Route) MakeDefaultPlan() *planner.Plan {
	plan := planner.NewPlan(planner.Plan{})

	// Value of 0.0 indicates that the route will never fail.
	plan.Percentage = fluent.NewMustFluentFloat("0.0")

	return &plan
}

func (route *Route) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		return planner.PLAN_SIGNAL_CONTINUE
	})

	return planner.CycleHooks{
		PreSleep: &preSleep,
	}
}

func (route *Route) Handle(w http.ResponseWriter, r *http.Request) {
	if route.Fault != nil {
		shouldSuccess := true

		for _, plan := range route.GetAssignedPlans() {
			if !*plan.GetCurrentValue().ComputedPercentageChance {
				shouldSuccess = false
				break
			}
		}

		if !shouldSuccess {
			route.Fault.Handle(w, r)
			return
		}
	}

	if route.Content.Whoami {
		w.Header().Set("Content-Type", "application/json")
		j := json.NewEncoder(w)
		j.SetIndent("", "  ")
		err := j.Encode(route.Content.GetReflectionContent(r))

		if err != nil {
			panic(err)
		}
		return
	}

	content := route.Content.Static

	if content == "" {
		content = "Hello from Kermoo!"
	}

	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(content))

	if err != nil {
		panic(err)
	}
}

func (route *Route) GetMethods() ([]string, error) {
	if len(route.Methods) == 0 {
		return []string{"HEAD", "GET", "POST"}, nil
	}

	validMethods := []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "CONNECT", "TRACE"}

	methods := []string{}

	for _, method := range route.Methods {
		method := strings.ToUpper(method)

		if !utils.Contains(validMethods, method) {
			return nil, fmt.Errorf("%s is not a valid HTTP method", method)
		}

		if !utils.Contains(methods, method) {
			methods = append(methods, method)
		}
	}

	return methods, nil
}

func (route *Route) Validate() error {
	if _, err := route.GetMethods(); err != nil {
		return err
	}

	if route.Fault != nil {
		if len(route.Fault.GetBadStatuses()) == 0 {
			return fmt.Errorf("route has no fault status - client and/or server errors needs to be enabled")
		}
	}

	return nil
}

type RouteContent struct {
	Static       string `json:"static"`
	Whoami       bool   `json:"whoami"`
	NoServerInfo bool   `json:"server_info"`
}

func (rc *RouteContent) GetReflectionContent(r *http.Request) ReflectorResponse {
	now := time.Now()

	server := ServerInfo{}

	if !rc.NoServerInfo {
		server = ServerInfo{
			Hostname:      os.Getenv("HOSTNAME"),
			InitializedAt: config.InitializedAt.Format(time.RFC3339Nano),
			CurrentTime:   now.Format(time.RFC3339Nano),
			UptimeSeconds: int64(now.Sub(config.InitializedAt).Seconds()),
			InterfaceIps:  utils.GetIpList(),
			KermooVersion: config.BuildVersion,
		}
	}

	return ReflectorResponse{
		Server: server,
		Request: RequestInfo{
			ConnectedFrom: r.RemoteAddr,
			Scheme:        r.URL.Scheme,
			Host:          r.Host,
			Path:          r.URL.Path,
			Query:         r.URL.Query(),
			Headers:       r.Header,
		},
	}
}

type ReflectorResponse struct {
	Server  ServerInfo  `json:"server"`
	Request RequestInfo `json:"request"`
}

type ServerInfo struct {
	Hostname      string   `json:"hostname"`
	InitializedAt string   `json:"initialized_at"`
	CurrentTime   string   `json:"current_time"`
	UptimeSeconds int64    `json:"uptime_seconds"`
	InterfaceIps  []string `json:"interface_ips"`
	KermooVersion string   `json:"kermoo_version"`
}

type RequestInfo struct {
	ConnectedFrom string              `json:"connected_from"`
	Scheme        string              `json:"scheme"`
	Host          string              `json:"host"`
	Path          string              `json:"path"`
	Query         map[string][]string `json:"query"`
	Headers       map[string][]string `json:"headers"`
}
