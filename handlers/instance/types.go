package instance

import "github.com/pivotal-pez/cfmgo"

type (
	//InstanceCreator - a type which can handle creating service instance requests
	InstanceCreator struct {
		Collection cfmgo.Collection
		Model      InstanceModel
	}
	InstanceModel struct {
		OrganizationGUID string                 `json:"organization_guid"`
		PlanID           string                 `json:"plan_id"`
		ServiceID        string                 `json:"service_id"`
		SpaceGUID        string                 `json:"space_guid"`
		Parameters       map[string]interface{} `json:"parameters"`
		TaskGUID         string
	}
)
