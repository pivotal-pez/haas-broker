package instance

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pivotal-pez/cfmgo"
)

const (
	//HandlerPath - path to normal instance handlers
	HandlerPath = "/v2/service_instances/{instance_id}"
	//AsyncHandlerPath - path to async poller
	AsyncHandlerPath = "/v2/service_instances/{instance_id}/last_operation"
)

//Put - handler function for put calls
func Put(collection cfmgo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	}
}

//Patch - handler function for patch calls
func Patch(collection cfmgo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	}
}

//Delete - handler function for delete calls
func Delete(collection cfmgo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	}
}

//Get - handler function for get calls
func Get(collection cfmgo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		category := vars["instance_id"]
		fmt.Fprintf(w, "Welcome to the home page!", category)
	}
}