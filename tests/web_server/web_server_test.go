package webserver_test

import (
	"buggybox/config"
	"buggybox/modules/web_server"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Testing routes
func TestRoutes(t *testing.T) {
	t.Run("test web server with static route", func(t *testing.T) {

		// Define the routes for the server
		routes := []web_server.Route{
			{
				Path:    "/info",
				Methods: []string{"GET"},
				Content: web_server.RouteContent{Static: "Hello, World!"},
			},
		}

		var (
			intf = "0.0.0.0"
			port = int32(8001)
		)

		ws := &web_server.WebServer{
			Routes:    routes,
			Interface: &intf,
			Port:      &port,
		}

		defer ws.Stop()

		ws.ListenOnBackground()

		// Give server a second to start
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://0.0.0.0:8001/info")
		if err != nil {
			t.Fatal(err)
		}

		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "Hello, World!", string(body))
	})

	t.Run("test web server with whoami route", func(t *testing.T) {
		// Define the routes for the server
		routes := []web_server.Route{
			{
				Path:    "/info",
				Methods: []string{"GET"},
				Content: web_server.RouteContent{Whoami: true},
			},
		}

		var (
			intf = "0.0.0.0"
			port = int32(8001)
		)

		ws := &web_server.WebServer{
			Routes:    routes,
			Interface: &intf,
			Port:      &port,
		}

		defer ws.Stop()

		ws.ListenOnBackground()

		// Give server a second to start
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://0.0.0.0:8001/info")
		if err != nil {
			t.Fatal(err)
		}

		body, _ := io.ReadAll(resp.Body)

		var response web_server.ReflectorResponse

		assert.NoError(t, json.Unmarshal([]byte(body), &response))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, config.BuildVersion, response.Server.BuggyboxVersion)
		assert.Equal(t, "/info", response.Request.Path)
	})
}
