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
	SchemaVersion string           `json:"schemaVersion"`
	Process       *process.Process `json:"process"`
	CpuLoad       *cpu.CpuLoader   `json:"cpuLoad"`
	MemoryLeak    *memory.MemoryLeak
	Plans         []*planner.Plan         `json:"plans"`
	WebServers    []*web_server.WebServer `json:"webServers"`
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
		if err := u.Process.Validate(); err != nil {
			return nil, fmt.Errorf("invalid process manager: %v", err)
		}

		prepared.Process = u.Process

		if err := prepared.preparePlannable(u.Process); err != nil {
			return nil, fmt.Errorf("unable to prepare process manager: %v", err)
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
