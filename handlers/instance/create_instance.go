package instance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pivotal-pez/haas-broker/handlers/catalog"
	"github.com/pivotal-pez/pezdispenser/pdclient"
)

//GetHandler - this is the handler that will be used for polling async
//provisioning status by the service broker
func (s *InstanceCreator) GetHandler(w http.ResponseWriter, req *http.Request) {
	responseBody := `{
		"state": "succeeded",
		"description": "Creating service (100% complete)."
	}`
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, responseBody)
}

//PutHandler - this is the actual handler method that will be used for the
//incoming request
func (s *InstanceCreator) PutHandler(w http.ResponseWriter, req *http.Request) {
	var (
		err          error
		bodyBytes    []byte
		statusCode   int
		responseBody string
	)
	s.parsePutVars(req)
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

func (s *InstanceCreator) parsePutVars(req *http.Request) {
	vars := mux.Vars(req)
	s.Model.InstanceID = vars["instance_id"]
}

func (s *InstanceCreator) getPlanName() string {
	return catalog.PlanGUIDMap[s.Model.PlanID]
}
