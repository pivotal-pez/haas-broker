package instance

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/haas-broker/handlers"
)

func setDashboardURL(vcapApp *cfenv.App) {
	if len(vcapApp.ApplicationURIs) > 0 {
		dashboardUrl = vcapApp.ApplicationURIs[0]
	}
}

//Put - handler function for put calls
func Put(collection cfmgo.Collection, dispenserCreds handlers.DispenserCreds, vcapApp *cfenv.App) func(http.ResponseWriter, *http.Request) {
	setDashboardURL(vcapApp)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	instanceCreator := &InstanceCreator{ClientDoer: client}
	instanceCreator.Collection = collection
	instanceCreator.Dispenser = dispenserCreds
	return instanceCreator.PutHandler
}

//Patch - handler function for patch calls
func Patch(collection cfmgo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	}
}

//Delete - handler function for delete calls
func Delete(collection cfmgo.Collection, dispenserCreds handlers.DispenserCreds) func(http.ResponseWriter, *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	instanceCreator := &InstanceCreator{ClientDoer: client}
	instanceCreator.Collection = collection
	instanceCreator.Dispenser = dispenserCreds
	return instanceCreator.DeleteHandler
}

//Get - handler function for get calls
func Get(collection cfmgo.Collection, dispenserCreds handlers.DispenserCreds) func(http.ResponseWriter, *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	instanceCreator := &InstanceCreator{ClientDoer: client}
	instanceCreator.Collection = collection
	instanceCreator.Dispenser = dispenserCreds
	return instanceCreator.GetHandler
}
