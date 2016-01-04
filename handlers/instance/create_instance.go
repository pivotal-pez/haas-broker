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

	var (
		err          error
		requestID    string
		responseBody = "{}"
		statusCode   int
	)
	s.Collection.Wake()
	s.parsePutVars(req)

	if requestID, err = GetRequestID(s.Model.InstanceID, s.Collection); err == nil {
		client := pdclient.NewClient(s.Dispenser.ApiKey, s.Dispenser.URL, s.ClientDoer)
		inventoryID := fmt.Sprintf("%s-%s", s.Model.OrganizationGUID, s.Model.SpaceGUID)
		meta := map[string]interface{}{
			RequestIDMetadataFieldname: requestID,
		}
		_, err = client.DeleteLease(s.Model.ServiceID, inventoryID, s.getPlanName(), meta)
		lo.G.Debug("deleteLease params: ", s.Model.ServiceID, inventoryID, s.getPlanName(), meta)
		statusCode = http.StatusOK
	}

	if err != nil {
		lo.G.Error("deletehandler error: ", err)
		statusCode = http.StatusGone
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, responseBody)
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

			if requestID, err := pdclient.GetRequestIDFromTaskResponse(task); err == nil {
				s.Model.UpdateField(s.Collection, RequestIDMetadataFieldname, requestID)
			}

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

			if leaseRes, _, err = client.PostLease(s.Model.ServiceID, inventoryID, s.getPlanName(), s.Model.OrganizationGUID, 14); err == nil {
				s.Model.TaskGUID = leaseRes.ID
				lo.G.Debug("leaserequest response: ", leaseRes)
				s.Model.RequestID, _ = pdclient.GetRequestIDFromTaskResponse(leaseRes)
				statusCode = http.StatusAccepted
				responseBody = fmt.Sprintf(`{"dashboard_url": "https://%s/show/%s"}`, dashboardUrl, s.Model.TaskGUID)
				s.Model.Save(s.Collection)
			}
		}
	}

	if err != nil {
		lo.G.Error("PutHandler error: ", err)
		statusCode = http.StatusNotAcceptable
		responseBody = fmt.Sprintf(`{"error_message": "%s"}`, err.Error())
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, responseBody)
}

func (s *InstanceCreator) parsePutVars(req *http.Request) {
	vars := mux.Vars(req)
	s.Model.InstanceID = vars[InstanceIDVarName]

	if v := req.FormValue("service_id"); v != "" {
		s.Model.ServiceID = v
	}

	if v := req.FormValue("plan_id"); v != "" {
		s.Model.PlanID = v
	}
}

func (s *InstanceCreator) getPlanName() string {
	return catalog.PlanGUIDMap[s.Model.PlanID]
}
