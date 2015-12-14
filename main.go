package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/nabeken/negroni-auth"
	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/haas-broker/handlers/binding"
	"github.com/pivotal-pez/haas-broker/handlers/catalog"
	"github.com/pivotal-pez/haas-broker/handlers/instance"
	"github.com/xchapter7x/lo"
)

func main() {
	lo.G.Debug("starting app")
	if appEnv, err := cfenv.Current(); err == nil {
		lo.G.Debug("parsed cfenv")
		serviceName := os.Getenv("MONGO_SERVICE_NAME")
		serviceURIName := os.Getenv("MONGO_SERVICE_URI_NAME")
		serviceURI := cfmgo.GetServiceBinding(serviceName, serviceURIName, appEnv)
		collectionName := os.Getenv("MONGO_COLLECTION_NAME")
		collection := cfmgo.Connect(cfmgo.NewCollectionDialer, serviceURI, collectionName)
		lo.G.Debug("created mongo conn", serviceURI, collectionName)
		n := negroni.Classic()
		lo.G.Debug("created negroni")

		if basicAuthService, err := appEnv.Services.WithName(os.Getenv("BASIC_AUTH_SERVICE_NAME")); err == nil {
			lo.G.Debug("parsed basic auth")
			user := basicAuthService.Credentials[os.Getenv("BASIC_AUTH_USERNAME_FIELD")].(string)
			pass := basicAuthService.Credentials[os.Getenv("BASIC_AUTH_PASSWORD_FIELD")].(string)
			n.Use(auth.Basic(user, pass))
			n.UseHandler(getRouter(collection))
			lo.G.Debug("starting server")
			n.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
			lo.G.Panic("run didnt lock!!!")

		} else {
			lo.G.Panic("error, could not find basic auth creds", err.Error())
		}

	} else {
		lo.G.Panic("error, failure to parse cfenv: ", err.Error())
	}
}

func getRouter(collection cfmgo.Collection) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(catalog.HandlerPath, catalog.Get()).Methods("GET")
	router.HandleFunc(instance.AsyncHandlerPath, instance.Get(collection)).Methods("GET")
	router.HandleFunc(instance.HandlerPath, instance.Put(collection)).Methods("PUT")
	router.HandleFunc(instance.HandlerPath, instance.Patch(collection)).Methods("PATCH")
	router.HandleFunc(instance.HandlerPath, instance.Delete(collection)).Methods("DELETE")
	router.HandleFunc(binding.HandlerPath, binding.Delete(collection)).Methods("DELETE")
	router.HandleFunc(binding.HandlerPath, binding.Put(collection)).Methods("PUT")
	return
}
