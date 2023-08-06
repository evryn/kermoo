package web_server

import (
	"buggybox/config"
	"buggybox/modules/planner"
	"buggybox/modules/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

type Route struct {
	planner.PlannableTrait
	Path     string        `json:"path"`
	Methods  []string      `json:"methods"`
	Content  RouteContent  `json:"content"`
	Plan     *planner.Plan `json:"plan"`
	PlanRefs []string      `json:"plan_refs"`
}

func (route *Route) GetUid() string {
	return slug.Make(fmt.Sprintf("route-%s", route.Path))
}

func (route *Route) GetDesiredPlanNames() []string {
	return route.PlanRefs
}

func (route *Route) HasCustomPlan() bool {
	return route.Plan != nil
}

func (route *Route) MakeCustomPlan() *planner.Plan {
	plan := *route.Plan
	return &plan
}

// Create a lifetime-long plan to serve route
func (route *Route) MakeDefaultPlan() *planner.Plan {
	// Value of 1.0 indicates that the route will always
	// be available.
	value := float32(1.0)

	plan := planner.InitPlan(planner.Plan{})
	plan.Value.Exactly = &value

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
	if route.Content.Whoami {
		w.Header().Set("Content-Type", "application/json")
		j := json.NewEncoder(w)
		j.SetIndent("", "  ")
		j.Encode(route.Content.GetReflectionContent(r))
		return
	}

	content := route.Content.Static

	if content == "" {
		content = "Hello from Kermoo!"
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
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
			Hostname: os.Getenv("HOSTNAME"),
			// TODO: InitializedAt:   time.InitialTime.Format(time.RFC3339Nano),
			CurrentTime: now.Format(time.RFC3339Nano),
			// TODO: UptimeSeconds:   int64(now.Sub(*Time.InitialTime).Seconds()),
			InterfaceIps:    utils.GetIpList(),
			BuggyboxVersion: config.BuildVersion,
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
	Hostname        string   `json:"hostname"`
	InitializedAt   string   `json:"initialized_at"`
	CurrentTime     string   `json:"current_time"`
	UptimeSeconds   int64    `json:"uptime_seconds"`
	InterfaceIps    []string `json:"interface_ips"`
	BuggyboxVersion string   `json:"buggybox_version"`
}

type RequestInfo struct {
	ConnectedFrom string              `json:"connected_from"`
	Scheme        string              `json:"scheme"`
	Host          string              `json:"host"`
	Path          string              `json:"path"`
	Query         map[string][]string `json:"query"`
	Headers       map[string][]string `json:"headers"`
}
