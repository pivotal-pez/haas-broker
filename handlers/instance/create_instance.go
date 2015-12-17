package instance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-pez/pezdispenser/pdclient"
)

//PutHandler - this is the actual handler method that will be used for the
//incoming request
func (s *InstanceCreator) PutHandler(w http.ResponseWriter, req *http.Request) {
	var (
		err          error
		bodyBytes    []byte
		statusCode   int
		responseBody string
	)

	if bodyBytes, err = ioutil.ReadAll(req.Body); err == nil {

		if err = json.Unmarshal(bodyBytes, &s.Model); err == nil {
			var (
				leaseRes pdclient.LeaseCreateResponseBody
			)
			client := pdclient.NewClient(s.Dispenser.ApiKey, s.Dispenser.URL, s.ClientDoer)
			inventoryID := fmt.Sprintf("%s-%s", s.Model.OrganizationGUID, s.Model.SpaceGUID)

			if leaseRes, _, err = client.PostLease(s.Model.ServiceID, inventoryID, s.getPlanName(), 14); err == nil {
				s.Model.TaskGUID = leaseRes.ID
				s.Model.Save(s.Collection)
				statusCode = http.StatusAccepted
				responseBody = fmt.Sprintf(`{"dashboard_url": "%s"}`, DashboardUrl)
			}
		}
	}

	if err != nil {
		statusCode = http.StatusNotAcceptable
		responseBody = fmt.Sprintf(`{"error_message": "%s"}`, err.Error())
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, responseBody)
}

func (s *InstanceCreator) getPlanName() string {
	return PlanGUIDMap[s.Model.PlanID]
}
