package web_server

import (
	"kermoo/modules/values"
	"math/rand"
	"net/http"
)

type RouteFault struct {
	PlanRefs []string `json:"planRefs"`

	Percentage values.MultiFloat `json:"percentage"`
	Interval   *values.Duration  `json:"interval"`
	Duration   *values.Duration  `json:"duration"`

	ResponseDelay values.Duration `json:"responseDelay"`
	ClientErrors  *bool           `json:"clientErrors"`
	ServerErrors  *bool           `json:"serverErrors"`
}

type RouteStatus struct {
	Code        int
	Description string
}

func (RouteFault *RouteFault) GetBadStatuses() []RouteStatus {
	statuses := []RouteStatus{}

	if RouteFault.ClientErrors != nil && *RouteFault.ClientErrors {
		statuses = append(statuses, []RouteStatus{
			{http.StatusBadRequest, "Bad Request"},
			{http.StatusUnauthorized, "Unauthorized"},
			{http.StatusForbidden, "Forbidden"},
			{http.StatusNotAcceptable, "Not Acceptable"},
			{http.StatusUnprocessableEntity, "Unprocessable Entity"},
			{http.StatusRequestTimeout, "Request Timeout"},
			{http.StatusConflict, "Conflict"},
		}...)
	}

	if RouteFault.ServerErrors == nil || *RouteFault.ServerErrors {
		statuses = append(statuses, []RouteStatus{
			{http.StatusInternalServerError, "Internal Server Error"},
			{http.StatusBadGateway, "Bad Gateway"},
			{http.StatusServiceUnavailable, "Service Unavailable"},
			{http.StatusGatewayTimeout, "Gateway Timeout"},
			{http.StatusInsufficientStorage, "Insufficient Storage"},
		}...)
	}

	return statuses
}

func (RouteFault *RouteFault) Handle(w http.ResponseWriter, r *http.Request) {
	statuses := RouteFault.GetBadStatuses()
	randomError := statuses[rand.Intn(len(statuses))]

	w.WriteHeader(randomError.Code)
	_, err := w.Write([]byte(randomError.Description))

	if err != nil {
		panic(err)
	}
}
