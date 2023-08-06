package e2e_test

import (
	"buggybox/modules/web_server"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebserverEndToEnd(t *testing.T) {
	t.Run("works with minimal config", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            schemaVersion: "0.1-beta"
            webServers:
            - port: 8080
              routes:
              - path: /my-world
		`, 2*time.Second)

		// Wait a few while for webserver to become available
		time.Sleep(500 * time.Millisecond)

		AssertHttpResponseContains(t, "GET", "http://0.0.0.0:8080/my-world", "Hello from Kermoo!")

		e2e.Wait()

		e2e.RequireTimedOut()
	})

	t.Run("works with more specific config", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            schemaVersion: "0.1-beta"
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
            schemaVersion: "0.1-beta"
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
            schemaVersion: "0.1-beta"
            webServers:
            - port: 8080
              routes:
              - path: /my-world
                content:
                  static: hello-world
              fault:
                plan:
                  interval: 100ms
                  value:
                    exactly: 0.5
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
            schemaVersion: "0.1-beta"
            plans:
            - name: disaster
              interval: 100ms
              value:
                exactly: 0.5
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
            schemaVersion: "0.1-beta"
            plans:
            - name: disaster
              subPlans:
              - value:
                  exactly: 0
                interval: 20ms
                duration: 1s
              - value:
                  between: [1, 1]
                interval: 20ms
                duration: 1s
              - value:
                  chart:
                    bars: [0.5, 0.4, 1]
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
}

func AssertHttpResponseContains(t *testing.T, method string, url string, expectedText string) {
	body, response, err := sendRequest(method, url, nil)
	require.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, body, expectedText)
}

func AssertHttpResponseCode(t *testing.T, method string, url string, expectedStatus int) {
	_, response, _ := sendRequest(method, url, nil)
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
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(respBody), resp, nil
}
