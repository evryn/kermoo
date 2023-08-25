package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kermoo/modules/web_server"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebserverEndToEnd(t *testing.T) {
	t.Run("works with defaults", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            webServers:
            - port: 8080
		`, 2*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		AssertHttpResponseContains(t, "GET", "http://0.0.0.0:8080/livez", "I'm Alive!")

		e2e.Wait()

		e2e.RequireTimedOut()
	})

	t.Run("works with more specific config", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            webServers:
            - port: 8000
              interface: 127.0.0.1
              routes:
              - path: /my-world
                methods: ["POST"]
                content:
                  static: hello-world
		`, 2*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		AssertHttpResponseCode(t, "GET", "http://127.0.0.1:8000/my-world", 405)
		AssertHttpResponseContains(t, "POST", "http://127.0.0.1:8000/my-world", "hello-world")

		e2e.Wait()

		e2e.RequireTimedOut()
	})

	t.Run("works with whoami responder", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.WithEnv("HOSTNAME=container-123")
		e2e.Start(`
            webServers:
            - port: 8000
              interface: 127.0.0.1
              routes:
              - path: /whoami
                content:
                  whoami: true
		`, 2*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		body, response, _ := sendRequest("GET", "http://127.0.0.1:8000/whoami", nil)

		assert.Equal(t, 200, response.StatusCode)

		whoami := web_server.ReflectorResponse{}
		require.NoError(t, json.Unmarshal([]byte(body), &whoami))
		assert.Equal(t, "container-123", whoami.Server.Hostname)
		assert.Equal(t, "Go-http-client/1.1", whoami.Request.Headers["User-Agent"][0])

		e2e.Wait()

		e2e.RequireTimedOut()
	})

	t.Run("works with a simple dedicated fault plan", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            webServers:
            - port: 8080
              routes:
              - path: /my-world
                content:
                  static: hello-world
              fault:
                interval: 100ms
                percentage: 50
		`, 3*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		success := GetWebserverSuccessRate(t, "GET", "http://0.0.0.0:8080/my-world", 20*time.Millisecond)

		e2e.Wait()

		assert.Greater(t, float32(0.9), success)
		assert.Less(t, float32(0.1), success)

		e2e.RequireTimedOut()
	})

	t.Run("works with a simple referenced fault plan", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            plans:
            - name: disaster
              interval: 100ms
              percentage: 50
            webServers:
            - port: 8080
              routes:
              - path: /my-world
                content:
                  static: hello-world
              fault:
                planRefs:
                - disaster
		`, 3*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		success := GetWebserverSuccessRate(t, "GET", "http://0.0.0.0:8080/my-world", 20*time.Millisecond)

		e2e.Wait()

		assert.Greater(t, float32(0.9), success)
		assert.Less(t, float32(0.1), success)

		e2e.RequireTimedOut()
	})

	t.Run("works with a complex referenced fault plan", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            plans:
            - name: disaster
              subPlans:
              - percentage: 100
                interval: 20ms
                duration: 1s
              - percentage: 0 to 0
                interval: 20ms
                duration: 1s
              - percentage: 50, 40, 100
                interval: 20ms
                duration: 1s
            webServers:
            - port: 8080
              routes:
              - path: /my-world
                content:
                  static: hello-world
              fault:
                planRefs:
                - disaster
		`, 5*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		phaseOneRate := GetWebserverSuccessRate(t, "GET", "http://0.0.0.0:8080/my-world", 1*time.Millisecond)

		time.Sleep(1000 * time.Millisecond)

		phaseTwoRate := GetWebserverSuccessRate(t, "GET", "http://0.0.0.0:8080/my-world", 1*time.Millisecond)

		time.Sleep(1000 * time.Millisecond)

		phaseThreeRate := GetWebserverSuccessRate(t, "GET", "http://0.0.0.0:8080/my-world", 5*time.Millisecond)

		e2e.Wait()

		assert.Equal(t, float32(0.0), phaseOneRate)
		assert.Equal(t, float32(1.0), phaseTwoRate)

		assert.Greater(t, float32(0.9), phaseThreeRate)
		assert.Less(t, float32(0.1), phaseThreeRate)

		e2e.RequireTimedOut()
	})

	t.Run("route fails with simple dedicated plan", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            webServers:
            - port: 8080
              routes:
              - path: /my-probe
                content:
                  static: hello-world
                fault:
                  interval: 100ms
                  percentage: 50
		`, 3*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		inspect := InspectRoute(t, "GET", "http://0.0.0.0:8080/my-probe", 20*time.Millisecond)

		e2e.Wait()

		assert.Equal(t, float32(0), inspect.WebserverErrorRate, "there must be no webserver error - its not intended to be faulty")
		assert.Equal(t, float32(0.0), inspect.ClientErrorRate, "there should be no client errors (4xx) since the default is disabled")
		assert.Less(t, float32(0.1), inspect.ServerErrorRate, "there should be at least some server errors (5xx)")
		assert.Greater(t, float32(0.9), inspect.SuccessRate, "there should be at least some failures")
		assert.Less(t, float32(0.1), inspect.SuccessRate, "there should be at least some success")

		e2e.RequireTimedOut()
	})

	t.Run("route fails with referenced plan and both client and server errors", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            plans:
            - name: readiness
              interval: 100ms
              percentage: 40
            webServers:
            - port: 8080
              routes:
              - path: /my-probe
                content:
                  static: hello-world
                fault:
                  planRefs:
                  - readiness
                  clientErrors: true
                  serverErrors: true
		`, 5*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		inspect := InspectRoute(t, "GET", "http://0.0.0.0:8080/my-probe", 40*time.Millisecond)

		e2e.Wait()

		assert.Equal(t, float32(0), inspect.WebserverErrorRate, "there must be no webserver error - its not intended to be faulty")
		assert.Less(t, float32(0.1), inspect.ClientErrorRate, "there should be at least some client errors (4xx)")
		assert.Less(t, float32(0.1), inspect.ServerErrorRate, "there should be at least some server errors (5xx)")
		assert.Greater(t, float32(0.9), inspect.SuccessRate, "there should be at least some failures")
		assert.Less(t, float32(0.1), inspect.SuccessRate, "there should be at least some success")

		e2e.RequireTimedOut()
	})
}

func AssertHttpResponseContains(t *testing.T, method string, url string, expectedText string) {
	body, response, err := sendRequest(method, url, nil)
	require.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, body, expectedText)
}

func AssertHttpResponseCode(t *testing.T, method string, url string, expectedStatus int) {
	_, response, err := sendRequest(method, url, nil)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus, response.StatusCode)
}

func GetWebserverSuccessRate(t *testing.T, method string, url string, sleep time.Duration) float32 {
	success := 0

	for i := 1; i <= 100; i++ {
		_, resp, _ := sendRequest(method, url, nil)

		if resp != nil {
			success++
			require.Equal(t, resp.StatusCode, 200, "route is faulty while it shouldn't be")
		}

		time.Sleep(sleep)
	}

	return float32(success) / float32(100)
}

type RouteInspection struct {
	SuccessRate        float32
	AverageDelay       time.Duration
	WebserverErrorRate float32
	ClientErrorRate    float32
	ServerErrorRate    float32
}

func InspectRoute(t *testing.T, method string, url string, sleep time.Duration) RouteInspection {
	total := 100
	totalDelay := time.Duration(0)
	success := 0
	clientErrors := 0
	serverErrors := 0
	webserverErrors := 0

	for i := 0; i < total; i++ {
		t1 := time.Now()
		_, resp, _ := sendRequest(method, url, nil)

		totalDelay += time.Since(t1)

		if resp == nil {
			webserverErrors++
		} else {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				clientErrors++
			} else if resp.StatusCode >= 500 && resp.StatusCode < 600 {
				serverErrors++
			} else {
				success++
			}
		}

		time.Sleep(sleep)
	}

	return RouteInspection{
		SuccessRate:        float32(success) / float32(total),
		AverageDelay:       time.Duration(totalDelay.Nanoseconds() / int64(total)),
		WebserverErrorRate: float32(webserverErrors) / float32(total),
		ClientErrorRate:    float32(clientErrors) / float32(total),
		ServerErrorRate:    float32(serverErrors) / float32(total),
	}
}

func sendRequest(method, url string, body []byte) (string, *http.Response, error) {
	// Create a new request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(respBody), resp, nil
}
