package router

type Response struct {
	Server  Server  `json:"server"`
	Request Request `json:"request"`
}

type Server struct {
	Hostname        string   `json:"hostname"`
	InitializedAt   string   `json:"initialized_at"`
	CurrentTime     string   `json:"current_time"`
	UptimeSeconds   int64    `json:"uptime_seconds"`
	InterfaceIps    []string `json:"interface_ips"`
	BuggyboxVersion string   `json:"buggybox_version"`
}

type Request struct {
	ConnectedFrom string              `json:"connected_from"`
	Scheme        string              `json:"scheme"`
	Host          string              `json:"host"`
	Path          string              `json:"path"`
	Query         map[string][]string `json:"query"`
	Headers       map[string][]string `json:"headers"`
}
