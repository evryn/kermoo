package user_config

import (
	"buggybox/modules/planner"
	"buggybox/modules/process"
	"buggybox/modules/web_server"
	"fmt"
)

type UserConfigType struct {
	SchemaVersion string                  `json:"schemaVersion"`
	Process       *process.Process        `json:"process"`
	Plans         []*planner.Plan         `json:"plans"`
	WebServers    []*web_server.WebServer `json:"webServers"`
}

func (u *UserConfigType) Validate() error {
	if u.SchemaVersion != "0.1-beta" {
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
			return nil, fmt.Errorf("unable to prepare manager: %v", err)
		}
	}

	for _, ws := range u.WebServers {
		// Prepare Web Server
		if err := ws.Validate(); err != nil {
			return nil, fmt.Errorf("invalid webserver %s: %v", ws.GetUid(), err)
		}

		prepared.WebServers = append(prepared.WebServers, ws)

		if err := prepared.preparePlannable(ws); err != nil {
			return nil, fmt.Errorf("unable to prepare webserver %s: %v", ws.GetUid(), err)
		}

		for _, route := range ws.Routes {
			// Prepare Routes
			if err := route.Validate(); err != nil {
				return nil, fmt.Errorf("invalid route %s for webserver %s: %v", route.GetUid(), ws.GetUid(), err)
			}

			if err := prepared.preparePlannable(route); err != nil {
				return nil, fmt.Errorf("unable to prepare route %s webserver %s: %v", route.GetUid(), ws.GetUid(), err)
			}
		}
	}

	return &prepared, nil
}
