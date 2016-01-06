package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	oauth2 "github.com/goincremental/negroni-oauth2"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/redisstore"
	"github.com/gorilla/mux"
	"github.com/nabeken/negroni-auth"
	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/haas-broker/handlers"
	"github.com/pivotal-pez/haas-broker/handlers/binding"
	"github.com/pivotal-pez/haas-broker/handlers/catalog"
	"github.com/pivotal-pez/haas-broker/handlers/instance"
	"github.com/unrolled/render"
	"github.com/xchapter7x/lo"
	"golang.org/x/net/context"
	baseOauth2 "golang.org/x/oauth2"
)

func main() {
	lo.G.Debug("starting app")

	if appEnv, err := cfenv.Current(); err == nil {
		lo.G.Debug("parsed cfenv")
		collection := getCollection(appEnv)
		dispenserCreds := getDispenserInfo(appEnv)
		n := negroni.Classic()
		setOauthContextInsecure(true)
		lo.G.Debug("created negroni")

		if oauthHandler, err := getOAuthHandler(); err == nil {

			if redisSession, err := getRedisSessionStore(appEnv); err == nil {
				n.Use(sessions.Sessions("haas_session", redisSession))

			} else {
				lo.G.Panic("could not create redis session: ", err)
			}
			n.Use(oauthHandler)
		} else {
			lo.G.Panic("not able to enable sso endpoints: ", err)
		}

		if router, err := getRouter(render.New(), collection, dispenserCreds, appEnv); err == nil {
			n.UseHandler(router)
			lo.G.Debug("starting server")
			n.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
			lo.G.Panic("run didnt lock!!!")

		} else {
			lo.G.Panic("error, could not properly create all routers and routes: ", err)
		}

	} else {
		lo.G.Panic("error, failure to parse cfenv: ", err.Error())
	}
}

func getRedisSessionStore(appEnv *cfenv.App) (rstore sessions.Store, err error) {
	var redisService *cfenv.Service

	if redisService, err = appEnv.Services.WithName(os.Getenv("REDIS_SERVICE_NAME")); err == nil {
		host := redisService.Credentials[os.Getenv("REDIS_HOST_FIELD")].(string)
		pass := redisService.Credentials[os.Getenv("REDIS_PASSWORD_FIELD")].(string)
		port := redisService.Credentials[os.Getenv("REDIS_PORT_FIELD")].(string)
		address := fmt.Sprintf("%s:%s", host, port)

		if rstore, err = redisstore.New(10, "tcp", address, pass, []byte("shhhhhhhdonttell")); err != nil {
			lo.G.Error("could not create a new redis store connection: ", err)
		}

	} else {
		lo.G.Error("couldnt find a redis service", err)
	}
	return
}

func getCollection(appEnv *cfenv.App) cfmgo.Collection {
	serviceName := os.Getenv("MONGO_SERVICE_NAME")
	serviceURIName := os.Getenv("MONGO_SERVICE_URI_NAME")
	serviceURI := cfmgo.GetServiceBinding(serviceName, serviceURIName, appEnv)
	collectionName := os.Getenv("MONGO_COLLECTION_NAME")
	lo.G.Debug("created mongo conn", serviceURI, collectionName)
	return cfmgo.Connect(cfmgo.NewCollectionDialer, serviceURI, collectionName)
}

func getDispenserInfo(appEnv *cfenv.App) handlers.DispenserCreds {
	serviceName := os.Getenv("DISPENSER_SERVICE_NAME")
	keyName := os.Getenv("DISPENSER_KEY_NAME")
	urlName := os.Getenv("DISPENSER_URL_NAME")
	service, _ := appEnv.Services.WithName(serviceName)
	creds := service.Credentials
	return handlers.DispenserCreds{
		ApiKey: creds[keyName].(string),
		URL:    creds[urlName].(string),
	}
}

func getBasicAuthCreds(appEnv *cfenv.App) (user, pass string, err error) {
	var basicAuthService *cfenv.Service

	if basicAuthService, err = appEnv.Services.WithName(os.Getenv("BASIC_AUTH_SERVICE_NAME")); err == nil {
		lo.G.Debug("parsed basic auth")
		user = basicAuthService.Credentials[os.Getenv("BASIC_AUTH_USERNAME_FIELD")].(string)
		pass = basicAuthService.Credentials[os.Getenv("BASIC_AUTH_PASSWORD_FIELD")].(string)

	} else {
		lo.G.Panic("error, could not find basic auth creds", err.Error())
	}
	return
}

func getRouter(renderer *render.Render, collection cfmgo.Collection, dispenserCreds handlers.DispenserCreds, appEnv *cfenv.App) (router *mux.Router, err error) {
	router = mux.NewRouter().StrictSlash(true)
	var (
		user string
		pass string
	)

	if user, pass, err = getBasicAuthCreds(appEnv); err == nil {
		v2Router := getV2Router(render.New(), collection, dispenserCreds, appEnv)
		router.PathPrefix(handlers.ServiceBrokerAPIVersion).Handler(negroni.New(
			negroni.HandlerFunc(auth.Basic(user, pass)),
			negroni.Wrap(v2Router),
		))

	} else {
		lo.G.Error("not enabling basic auth endpoints: ", err)
	}
	ssoRouter := getSSORouter(render.New(), collection, dispenserCreds)

	router.PathPrefix(instance.SSOPathPrefix).Handler(negroni.New(
		oauth2.LoginRequired(),
		negroni.Wrap(ssoRouter),
	))
	return
}

func getSSORouter(renderer *render.Render, collection cfmgo.Collection, dispenserCreds handlers.DispenserCreds) (ssoRouter *mux.Router) {
	ssoRouter = mux.NewRouter().PathPrefix(instance.SSOPathPrefix).Subrouter().StrictSlash(true)
	ssoRouter.HandleFunc(instance.ServiceInstanceDash, instance.GetDashboard(dispenserCreds, collection, render.New())).Methods("GET")
	ssoRouter.HandleFunc("/oauth2callback", instance.GetDashboard(dispenserCreds, collection, render.New())).Methods("GET")
	return
}

func getOAuthHandler() (uaaProvider negroni.Handler, err error) {
	var (
		app          *cfenv.App
		oauthService *cfenv.Service
	)

	if app, err = cfenv.Current(); err == nil {
		services := app.Services

		if oauthService, err = services.WithName(os.Getenv("OAUTH_SERVICE_NAME")); err == nil {
			clientID := oauthService.Credentials[os.Getenv("OAUTH_CLIENT_FIELD")].(string)
			clientSecret := oauthService.Credentials[os.Getenv("OAUTH_CLIENT_SECRET_FIELD")].(string)
			authzEndpoint := oauthService.Credentials[os.Getenv("OAUTH_AUTHZ_ENDPOINT_FIELD")].(string)
			tokenEndpoint := oauthService.Credentials[os.Getenv("OAUTH_TOKEN_ENDPOINT_FIELD")].(string)

			oauthOpts := &oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Scopes:       []string{"cloud_controller_service_permissions.read", "openid"},
			}
			uaaProvider = negroni.HandlerFunc(oauth2.NewOAuth2Provider(oauthOpts, authzEndpoint, tokenEndpoint))
		} else {
			lo.G.Error("could not find oauthservice: ", err)
		}
	} else {
		lo.G.Error("could not grab valid vcap", err)
	}
	return
}

func setOauthContextInsecure(ignore bool) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ignore},
	}
	client := &http.Client{Transport: tr}
	baseOauth2.NoContext = context.WithValue(context.Background(), baseOauth2.HTTPClient, client)
}

func getV2Router(renderer *render.Render, collection cfmgo.Collection, dispenserCreds handlers.DispenserCreds, appEnv *cfenv.App) (v2Router *mux.Router) {
	v2Router = mux.NewRouter().PathPrefix(handlers.ServiceBrokerAPIVersion).Subrouter().StrictSlash(true)
	v2Router.HandleFunc(catalog.HandlerPath, catalog.Get()).Methods("GET")
	v2Router.HandleFunc(instance.AsyncHandlerPath, instance.Get(collection, dispenserCreds)).Methods("GET")
	v2Router.HandleFunc(instance.HandlerPath, instance.Put(collection, dispenserCreds, appEnv)).Methods("PUT")
	v2Router.HandleFunc(instance.HandlerPath, instance.Patch(collection)).Methods("PATCH")
	v2Router.HandleFunc(instance.HandlerPath, instance.Delete(collection, dispenserCreds)).Methods("DELETE")
	v2Router.HandleFunc(binding.HandlerPath, binding.Delete(collection)).Methods("DELETE")
	v2Router.HandleFunc(binding.HandlerPath, binding.Put(collection)).Methods("PUT")
	return
}
