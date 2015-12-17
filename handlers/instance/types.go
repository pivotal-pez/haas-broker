package instance

import (
	"net/http"

	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/haas-broker/handlers"
)

type (
	//InstanceCreator - a type which can handle creating service instance requests
	InstanceCreator struct {
		Collection cfmgo.Collection
		Dispenser  handlers.DispenserCreds
		ClientDoer clientDoer
		Model      InstanceModel
	}
	//InstanceModel - persistence model struct
	InstanceModel struct {
		OrganizationGUID string                 `json:"organization_guid"`
		PlanID           string                 `json:"plan_id"`
		ServiceID        string                 `json:"service_id"`
		SpaceGUID        string                 `json:"space_guid"`
		Parameters       map[string]interface{} `json:"parameters"`
		TaskGUID         string
	}
	clientDoer interface {
		Do(req *http.Request) (resp *http.Response, err error)
	}
)
