package catalog

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
)

const (
	//HandlerPath - path for catalog handler to register against
	HandlerPath = "/catalog"
)

func getDashboardInfo() (clientID string, clientSecret string, redirectURL string) {

	var (
		app          *cfenv.App
		oauthService *cfenv.Service
		err          error
	)

	if app, err = cfenv.Current(); err == nil {
		redirectURL = app.ApplicationURIs[0]
		services := app.Services

		if oauthService, err = services.WithName(os.Getenv("OAUTH_SERVICE_NAME")); err == nil {
			clientID = oauthService.Credentials[os.Getenv("OAUTH_CLIENT_FIELD")].(string)
			clientSecret = oauthService.Credentials[os.Getenv("OAUTH_CLIENT_SECRET_FIELD")].(string)
		}
	}
	return
}

//Get - function to handle a get request
func Get() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		clientID, clientSecret, redirectURL := getDashboardInfo()
		response := fmt.Sprintf(`{
			"services": [{
				"id": "5a9b9f22-a08d-11e5-8062-7831c1d4f660",
				"name": "pez-haas",
				"description": "Lease on-demand hardware as a service",
				"metadata":{
            "displayName":"PEZ-HaaS",
            "imageUrl":"http://s12.postimg.org/wt91ic9pp/broker_icon.png",
						"providerDisplayName":"PEZ"
         },
				"bindable": false,
				"plans": [{
					"id": "6a977311-a08d-11e5-8062-7831c1d4f660",
					"name": "m1.small",
					"description": "A small instance of hardware as a service",
					"metadata":{
						"bullets":[
							"96gb memory (min)",
							"2.7 GHz x (4 sockets / 12 cores per)",
							"3TB NFS shared storage",
							"40 TB total local disk",
							"/24 network (on 10.65.x.x pivotal vpn)",
							"ESXi installed"
						]
					}
				}],
				"dashboard_client": {
					"id": "%s",
					"secret": "%s",
					"redirect_uri": "https://%s/oauth2callback"
				}
			}]
		}`, clientID, clientSecret, redirectURL)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, response)
	}
}
