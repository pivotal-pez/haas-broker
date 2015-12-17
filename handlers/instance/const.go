package instance

const (
	//HandlerPath - path to normal instance handlers
	HandlerPath = "/v2/service_instances/{instance_id}"
	//AsyncHandlerPath - path to async poller
	AsyncHandlerPath = "/v2/service_instances/{instance_id}/last_operation"
)

var (
	dashboardUrl = "https://www.pezapp.io"
)
