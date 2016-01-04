package binding

import (
	"fmt"
	"net/http"

	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/haas-broker/handlers"
)

var (
	//HandlerPath - path to normal instance handlers
	HandlerPath = fmt.Sprintf("%s/service_instances/{instance_id}/service_bindings/{binding_id}", handlers.ServiceBrokerAPIVersion)
)

//Put - handler function for put calls
func Put(collection cfmgo.Collection) func(http.ResponseWriter, *http.Request) {
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
