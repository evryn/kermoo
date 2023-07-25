package Router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func MustSetupRouter(addr string) {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/readyz", readyzHandler)
	r.HandleFunc("/livez", livezHandler)
	http.ListenAndServe(":8080", r)
}

type Response struct {
	Hostname string
	Headers  map[string][]string `json:"headers"`
	Path     string              `json:"path"`
	Query    map[string][]string `json:"query"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Hostname: os.Getenv("HOSTNAME"),
		Headers:  r.Header,
		Path:     r.URL.Path,
		Query:    r.URL.Query(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Println(w, "HTTP: GET `/` 200 OK")
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func livezHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
