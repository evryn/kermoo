package webserver_test

import (
	"kermoo/modules/utils"
	"kermoo/modules/web_server"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of planner.PlannableTrait for testing purposes
type MockPlannableTrait struct{}

func (m MockPlannableTrait) Plan() error {
	return nil
}

// Mock implementation of RouteFault for testing purposes
type MockRouteFault struct{}

func (f *MockRouteFault) GetBadStatuses() []int {
	return []int{400, 500}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		route   web_server.Route
		wantErr bool
	}{
		{
			name: "valid defaults",
			route: web_server.Route{
				Path:    "/api",
				Methods: nil,
				Content: web_server.RouteContent{},
				Fault:   nil,
			},
			wantErr: false,
		},
		{
			name: "valid with specific methods",
			route: web_server.Route{
				Path:    "/api",
				Methods: []string{"GET", "POST"},
				Content: web_server.RouteContent{},
				Fault:   nil,
			},
			wantErr: false,
		},
		{
			name: "valid with specified default fault",
			route: web_server.Route{
				Path:    "/api",
				Methods: []string{"GET", "POST"},
				Content: web_server.RouteContent{},
				Fault:   &web_server.RouteFault{},
			},
			wantErr: false,
		},
		{
			name: "invalid with no fault status",
			route: web_server.Route{
				Path:    "/api",
				Methods: []string{"GET", "POST"},
				Content: web_server.RouteContent{},
				Fault: &web_server.RouteFault{
					ClientErrors: utils.NewP[bool](false),
					ServerErrors: utils.NewP[bool](false),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid with wrong methods",
			route: web_server.Route{
				Path:    "/api",
				Methods: []string{"AAAA"},
				Content: web_server.RouteContent{},
				Fault:   nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Error(t, tt.route.Validate())
			} else {
				assert.Nil(t, tt.route.Validate())
			}
		})
	}
}
