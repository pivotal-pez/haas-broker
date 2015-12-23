package instance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xchapter7x/lo"

	"github.com/gorilla/mux"
	"github.com/pivotal-pez/haas-broker/handlers/catalog"
	"github.com/pivotal-pez/pezdispenser/pdclient"
)

//DeleteHandler - this is the handler that will be used when a user deletes a
//service instance
func (s *InstanceCreator) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	s.Collection.Wake()
	s.parsePutVars(req)

}

//GetHandler - this is the handler that will be used for polling async
//provisioning status by the service broker
func (s *InstanceCreator) GetHandler(w http.ResponseWriter, req *http.Request) {
	var (
		err          error
		taskID       string
		responseBody string
		task         pdclient.TaskResponse
	)
	s.Collection.Wake()
	s.parsePutVars(req)

	if taskID, err = GetTaskID(s.Model.InstanceID, s.Collection); err == nil {
		client := pdclient.NewClient(s.Dispenser.ApiKey, s.Dispenser.URL, s.ClientDoer)

		if task, _, err = client.GetTask(taskID); err == nil {

			switch task.Status {
			case TaskStatusComplete:
				responseBody = fmt.Sprintf(SuccessGetHandlerBody, task.Status)
			case TaskStatusFailed:
				responseBody = fmt.Sprintf(FailureGetHandlerBody, task.Status)
			default:
				responseBody = fmt.Sprintf(PendingGethandlerBody, task.Status)
			}
		}
	}
	if err != nil {
		lo.G.Error("gethandler error: ", err)
		responseBody = fmt.Sprintf(PendingGethandlerBody, err.Error())
	}
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
	s.Collection.Wake()
	s.parsePutVars(req)
	if bodyBytes, err = ioutil.ReadAll(req.Body); err == nil {

		if err = json.Unmarshal(bodyBytes, &s.Model); err == nil {
			var (
				leaseRes pdclient.TaskResponse
			)
			client := pdclient.NewClient(s.Dispenser.ApiKey, s.Dispenser.URL, s.ClientDoer)
			inventoryID := fmt.Sprintf("%s-%s", s.Model.OrganizationGUID, s.Model.SpaceGUID)

			if leaseRes, _, err = client.PostLease(s.Model.ServiceID, inventoryID, s.getPlanName(), 14); err == nil {
				s.Model.TaskGUID = leaseRes.ID
				s.Model.RequestID = pdclient.GetRequestIDFromTaskResponse(leaseRes)
				s.Model.Save(s.Collection)
				statusCode = http.StatusAccepted
				responseBody = fmt.Sprintf(`{"dashboard_url": "https://%s/show/%s"}`, dashboardUrl, s.Model.TaskGUID)
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
	s.Model.InstanceID = vars[InstanceIDVarName]
}

func (s *InstanceCreator) getPlanName() string {
	return catalog.PlanGUIDMap[s.Model.PlanID]
}
