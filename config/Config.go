package config

import "time"

type PlannerDefault struct {
	Minimum  float32
	Maximum  float32
	Interval time.Duration
}
type DefaultTemplate struct {
	Planner PlannerDefault
}

var (
	AppTitle       = "📦🐛 BuggyBox 🐛📦"
	AppDescription = "An app with the purpose of demonstating real-world malfunctioning applications. Good for testing and learnign container management topics."
	BuildVersion   string
	BuildRef       string
	BuildDate      string
	Default        = DefaultTemplate{
		Planner: PlannerDefault{
			Minimum:  0.0,
			Maximum:  1.0,
			Interval: 2 * time.Second,
		},
	}
)
