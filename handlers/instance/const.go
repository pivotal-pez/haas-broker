package instance

import "errors"

const (
	//HandlerPath - path to normal instance handlers
	HandlerPath = "/v2/service_instances/{instance_id}"
	//AsyncHandlerPath - path to async poller
	AsyncHandlerPath = "/v2/service_instances/{instance_id}/last_operation"
	//TaskStatusComplete ---
	TaskStatusComplete = "complete"
	//AgentTaskStatusFailed ---
	TaskStatusFailed = "failed"
	//CollectionInstanceIDQueryField --
	CollectionInstanceIDQueryField = "instanceid"
	//SuccessGetHandlerBody --
	SuccessGetHandlerBody = `{
		"state": "succeeded",
		"description": "%s"
	}`
	//FailureGetHandlerBody --
	FailureGetHandlerBody = `{
		"state": "failed",
		"description": "%s"
	}`
	//PendingGethandlerBody --
	PendingGethandlerBody = `{
		"state": "in progress",
		"description": "%s"
	}`
)

var (
	dashboardUrl         = "https://www.pezapp.io"
	ErrInvalidInstanceID = errors.New("invalid instance id while attempting to get taskid")
)
