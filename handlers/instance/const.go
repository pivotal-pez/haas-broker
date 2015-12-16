package instance

import "net/http"

const (
	//HandlerPath - path to normal instance handlers
	HandlerPath = "/v2/service_instances/{instance_id}"
	//AsyncHandlerPath - path to async poller
	AsyncHandlerPath = "/v2/service_instances/{instance_id}/last_operation"
	//DashboardUrl - the url to the service instance tracker dashboard
	DashboardUrl = "https://www.pezapp.io"
)

var (
	HttpClient  clientDoer = new(http.Client)
	PlanGUIDMap            = map[string]string{
		"6a977311-a08d-11e5-8062-7831c1d4f660": "m1.small",
	}
)
