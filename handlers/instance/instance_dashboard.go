package instance

import (
	"crypto/tls"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/haas-broker/handlers"
	"github.com/pivotal-pez/pezdispenser/pdclient"
	"github.com/unrolled/render"
)

//GetDashboard - get handler for dashboard with a http.Client
func GetDashboard(dispenserCreds handlers.DispenserCreds, collection cfmgo.Collection, renderer *render.Render) func(http.ResponseWriter, *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpDoer := &http.Client{Transport: tr}
	return BuildDashboardHandler(dispenserCreds, collection, renderer, httpDoer)
}

//BuildDashboardHandler - build an initialized get handler for the dashboard
func BuildDashboardHandler(dispenserCreds handlers.DispenserCreds, collection cfmgo.Collection, renderer *render.Render, httpDoer clientDoer) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		collection.Wake()
		vars := mux.Vars(req)
		taskGUID := vars[TaskGUIDVarName]
		client := pdclient.NewClient(dispenserCreds.ApiKey, dispenserCreds.URL, httpDoer)
		response := PendingResponse

		if task, _, err := client.GetTask(taskGUID); err == nil {
			switch task.Status {
			case TaskStatusFailed:
				response = FailureResponse

			case TaskStatusComplete:
				if status, ok := task.MetaData["status"]; ok {
					if data, ok := status.(map[string]interface{})["data"]; ok {
						responseBytes, _ := json.MarshalIndent(data, "", "  ")
						response = string(responseBytes[:])
					}
				}
			}
		}
		renderer.HTML(res, http.StatusOK, "haas", response)
	}
}
