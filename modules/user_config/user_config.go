package user_config

import (
	"fmt"
	"kermoo/modules/cpu"
	"kermoo/modules/memory"
	"kermoo/modules/planner"
	"kermoo/modules/process"
	"kermoo/modules/web_server"
)

type UserConfigType struct {
	// SchemaVersion is an optional indication of the schema version of the current
	// given configuration. It's here for future backwards compatibility.
	SchemaVersion string `json:"schemaVersion"`

	// Process optionally controls the execution of the main process. It's to determine an initial
	// delay and/or sudden exit of the process in the given time.
	//
	// By default, the process won't have any additional delay or perform sudden exit.
	Process *process.Process `json:"process"`

	// CpuLoad optionally simulates the CPU load of the machine. You can specify interval, duration
	// and load of it in percentage.
	//
	// By default, no CPU load is simulated.
	CpuLoad *cpu.CpuLoader `json:"cpuLoad"`

	// MemoryLeak optionally simulates the memory leak by consuming the memory of the machine.
	// You can specify interval, duration and size of the leak.
	//
	// By default, no memory leak is simulated.
	MemoryLeak *memory.MemoryLeak

	// WebServers is an optional array of web servers that will be used to serve defined routes.
	// It can be configured to fail with percentage over an specific duration of time with specific
	// interval. Routes can be configured too.
	//
	// By default, no web server is initiated.
	WebServers []*web_server.WebServer `json:"webServers"`

	// Plans is an optional array of plans which is there to avoid re-defining some repeatitive
	// failure plans. It can be refered from a webServer, route, cpuLoad, or memoryLeak.
	Plans []*planner.Plan `json:"plans"`
}

func (u *UserConfigType) Validate() error {
	if u.SchemaVersion != "" && u.SchemaVersion != "1" {
		return fmt.Errorf("schema version is not supported")
	}

	return nil
}

func (u *UserConfigType) GetPreparedConfig() (*PreparedConfigType, error) {
	prepared := PreparedConfigType{
		Plans: u.Plans,
	}

	// Prepare process manager
	if u.Process != nil {
		prepared.Process = u.Process

		if u.Process.Exit != nil {
			if err := u.Process.Validate(); err != nil {
				return nil, fmt.Errorf("invalid process manager: %v", err)
			}

			if err := prepared.preparePlannable(u.Process); err != nil {
				return nil, fmt.Errorf("unable to prepare process manager: %v", err)
			}
		}

	}

	// Prepare CPU Load
	if u.CpuLoad != nil {
		if err := u.CpuLoad.Validate(); err != nil {
			return nil, fmt.Errorf("invalid cpu load: %v", err)
		}

		if err := prepared.preparePlannable(u.CpuLoad); err != nil {
			return nil, fmt.Errorf("unable to prepare cpu load: %v", err)
		}
	}

	// Prepare Memory Leaker
	if u.MemoryLeak != nil {
		if err := u.MemoryLeak.Validate(); err != nil {
			return nil, fmt.Errorf("invalid memory leaker: %v", err)
		}

		if err := prepared.preparePlannable(u.MemoryLeak); err != nil {
			return nil, fmt.Errorf("unable to prepare memory leaker: %v", err)
		}
	}

	// Prepare Web Server
	if err := u.prepareWebservers(&prepared); err != nil {
		return nil, err
	}

	return &prepared, nil
}

func (u *UserConfigType) prepareWebservers(p *PreparedConfigType) error {
	for _, ws := range u.WebServers {
		// Prepare Web Server
		if err := ws.Validate(); err != nil {
			return fmt.Errorf("invalid webserver %s: %v", ws.GetName(), err)
		}

		p.WebServers = append(p.WebServers, ws)

		if err := p.preparePlannable(ws); err != nil {
			return fmt.Errorf("unable to prepare webserver %s: %v", ws.GetName(), err)
		}

		for _, route := range ws.Routes {
			// Prepare Routes
			if err := route.Validate(); err != nil {
				return fmt.Errorf("invalid route %s for webserver %s: %v", route.GetName(), ws.GetName(), err)
			}

			if err := p.preparePlannable(route); err != nil {
				return fmt.Errorf("unable to prepare route %s webserver %s: %v", route.GetName(), ws.GetName(), err)
			}
		}
	}

	return nil
}
