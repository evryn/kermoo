package web_server

import (
	"buggybox/config"
	"buggybox/modules/utils"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type Route struct {
	Path    string       `json:"path"`
	Methods []string     `json:"methods"`
	Content RouteContent `json:"content"`
}

func (route *Route) Handle(w http.ResponseWriter, r *http.Request) {
	if route.Content.Reflector {
		w.Header().Set("Content-Type", "application/json")
		j := json.NewEncoder(w)
		j.SetIndent("", "  ")
		j.Encode(route.Content.GetReflectionContent(r))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(route.Content.Static))
}

type RouteContent struct {
	Static       string `json:"static"`
	Reflector    bool   `json:"reflector"`
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
