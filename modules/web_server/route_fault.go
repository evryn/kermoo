package web_server

import (
	"kermoo/modules/fluent"
	"math/rand"
	"net/http"
)

type RouteFault struct {
	// PlanRefs is an optional list of plan names. It can used to avoid redundant
	// re-declearing of plans in large-scale configurations.
	// PlanRefs overrides Size, Interval and Duration fields are overrided in favor
	// of the one defined in the referenced plan.
	//
	// When more than one plan is referenced, the route will use AND operator to
	// determine if it should fail or not. This is useful when for example your
	// `/readyz` route wants to have its own failure plan but explicitly depends on
	// another plan such as the one for `/livez`
	PlanRefs []string `json:"planRefs"`

	// Percentage determines the chance of failing. 0 means no not failure at all and 100
	// means always failing. By failing, we mean the route will respond with either 4xx (client
	// side error) or 5xx (server side error) - which is also configurable.
	//
	// For specific and ranged declearations, it's going to use that but when an array of
	// percentages are specified, it'll act like a graph of bars and iterate over them.
	Percentage fluent.FluentFloat `json:"percentage"`

	// Interval decides how long each desicion to stay failing or serving should last.
	// A value above one second is recommended but you're free  to use any interval.
	//
	// Default is one second.
	Interval *fluent.FluentDuration `json:"interval"`

	// Duration defines the duration of the entire route. Leave it empty for
	// life-long running or specify one to end the route completely after that and last decision
	// will be happening for ever.
	// In fact, Duration/Interval determines the number of cycle, if defined. Default is empty
	// for unlimited activity.
	Duration *fluent.FluentDuration `json:"duration"`

	// ResponseDelay adds a delay to each response - no matter if its in good or bad state.
	ResponseDelay fluent.FluentDuration `json:"responseDelay"`

	// ClientErrors indicates when the route is in failing state, it can respond with 4xx (client
	// side errors).
	//
	// Default is false so that only 5xx errors are returned.
	ClientErrors *bool `json:"clientErrors"`

	// ServerErrors indicates when the route is in failing state, it can respond with 5xx (server
	// side errors).
	//
	// Default is true.
	ServerErrors *bool `json:"serverErrors"`
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
