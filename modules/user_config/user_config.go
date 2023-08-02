package user_config

import (
	"buggybox/modules/planner"
	"buggybox/modules/process"
	"buggybox/modules/web_server"
	"fmt"
)

var UserConfig UserConfigType

type UserConfigType struct {
	SchemaVersion string                 `json:"schemaVersion"`
	Process       process.Process        `json:"process"`
	Plans         []planner.Plan         `json:"plans"`
	WebServers    []web_server.WebServer `json:"webServers"`
}

func (u *UserConfigType) Validate() error {
	if u.SchemaVersion != "0.1-beta" {
		return fmt.Errorf("schema version is not supported")
	}

	return nil
}
