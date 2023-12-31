package config

import (
	"time"
)

type PlannerDefault struct {
	Minimum  float32
	Maximum  float32
	Interval time.Duration
}

type WebServerDefault struct {
	Interface string
	Port      int32
}

type DefaultTemplate struct {
	Planner   PlannerDefault
	WebServer WebServerDefault
}

var (
	AppTitle       = "📦🐛 Kermoo 🐛📦"
	AppDescription = "An app with the purpose of demonstating real-world malfunctioning applications. Good for testing and learnign container management topics."
	BuildVersion   string
	BuildRef       string
	BuildDate      string
	Default        = DefaultTemplate{
		Planner: PlannerDefault{
			Interval: 1 * time.Second,
		},
		WebServer: WebServerDefault{
			Port:      80,
			Interface: "0.0.0.0",
		},
	}
)

var InitializedAt time.Time
