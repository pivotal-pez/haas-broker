package catalog

import (
	"fmt"

	"github.com/pivotal-pez/haas-broker/handlers"
)

var (
	M1SmallPlanID = "6a977311-a08d-11e5-8062-7831c1d4f660"
	PlanGUIDMap   = map[string]string{
		M1SmallPlanID: "m1.small",
	}
	//HandlerPath - path for catalog handler to register against
	HandlerPath = fmt.Sprintf("%s/catalog", handlers.ServiceBrokerAPIVersion)
)
