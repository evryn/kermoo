package Router

import (
	"buggybox/config"
	"buggybox/modules/Time"
	"buggybox/modules/Utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func MustSetupRouter(addr string) {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/readyz", readyzHandler)
	r.HandleFunc("/livez", livezHandler)
	fmt.Printf("Starting HTTP server on %s\n", addr)
	panic(http.ListenAndServe(addr, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	response := Response{
		Server: Server{
			Hostname:        os.Getenv("HOSTNAME"),
			InitializedAt:   Time.InitialTime.Format(time.RFC3339Nano),
			CurrentTime:     now.Format(time.RFC3339Nano),
			UptimeSeconds:   int64(now.Sub(*Time.InitialTime).Seconds()),
			InterfaceIps:    Utils.GetIpList(),
			BuggyboxVersion: config.BuildVersion,
		},
		Request: Request{
			ConnectedFrom: r.RemoteAddr,
			Scheme:        r.URL.Scheme,
			Host:          r.Host,
			Path:          r.URL.Path,
			Query:         r.URL.Query(),
			Headers:       r.Header,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	j.SetIndent("", "  ")
	j.Encode(response)

	fmt.Println("HTTP: `/` 200 OK")
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func livezHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
