// Copyright (c) 2021 Acronis International GmbH
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// The `Server` package initiliaze routes and run the server
package server

//nolint
import (
	"fmt"
	"log"
	"net/http"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/server/handlers"
	"github.com/gorilla/mux"
)

// The `App` struct is used to hold application level details such as request routes
type App struct {
	Router *mux.Router
}

// Initialize routes with matched http request paths
func (app *App) InitializeRoutes(appConfig *config.Config) error {
	app.Router = mux.NewRouter()
	return app.initializeRoutes(appConfig)
}

// Listen and run the service on particular port
func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

// Initialize routes for Tenants, Users, OIs, AccessPolicies and Usage Endpoints
func (app *App) initializeRoutes(appConfig *config.Config) error {
	app.initializeTenantRoutes()
	app.initializeUserRoutes()
	app.initializeOfferingItemsRoutes()
	app.initializeUsageRoutes()
	app.initializeAccessPolicyRoutes()

	if appConfig.AuthSettings.Enabled {
		// Auth routes
		authHandler, err := handlers.NewAuthHandler(appConfig)
		if err != nil {
			return fmt.Errorf("error creating auth handler. %w", err)
		}
		app.Router.HandleFunc("/auth/acronis/redirect", authHandler.AuthInitiateRedirect).Methods(http.MethodGet)
		app.Router.HandleFunc("/auth/acronis/callback", authHandler.AuthCallback).Methods(http.MethodGet)
		app.Router.HandleFunc("/auth/acronis/initiate", authHandler.AuthInitiate).Methods(http.MethodGet)
		app.Router.HandleFunc("/auth/logout", authHandler.Logout).Methods(http.MethodPost)

		app.initializeWebRoutes(authHandler)
	}

	return nil
}

func (app *App) initializeUserRoutes() {
	app.Router.HandleFunc("/users", handlers.FetchUsers).Methods(http.MethodGet)
	app.Router.HandleFunc("/users", handlers.CreateOrUpdateUser).Methods(http.MethodPost)
	app.Router.HandleFunc("/users/{id}", handlers.FetchUser).Methods(http.MethodGet)
	app.Router.HandleFunc("/users/{id}", handlers.DeleteUser).Methods(http.MethodDelete)
}

func (app *App) initializeOfferingItemsRoutes() {
	app.Router.HandleFunc("/offering_items", handlers.GetOfferingItems).Methods("GET")
	app.Router.HandleFunc("/tenants/{tenant_id}/offering_items", handlers.CreateOrUpdateOfferingItem).Methods("POST")
	app.Router.HandleFunc("/tenants/{tenant_id}/offering_items/{offering_item_name}", handlers.GetOfferingItem).Methods("GET")
	app.Router.HandleFunc("/tenants/{tenant_id}/offering_items/{offering_item_name}", handlers.DeleteOfferingItem).Methods("DELETE")
}

func (app *App) initializeTenantRoutes() {
	app.Router.HandleFunc("/tenants", handlers.FetchTenants).Methods(http.MethodGet)
	app.Router.HandleFunc("/tenants/{tenant_id}", handlers.FetchTenant).Methods(http.MethodGet)
	app.Router.HandleFunc("/tenants", handlers.CreateOrUpdateTenant).Methods(http.MethodPost)
	app.Router.HandleFunc("/tenants/{tenant_id}", handlers.DeleteTenant).Methods(http.MethodDelete)
}

func (app *App) initializeAccessPolicyRoutes() {
	app.Router.HandleFunc("/access_policies", handlers.FetchAccessPolicies).Methods(http.MethodGet)
	app.Router.HandleFunc("/access_policies/{id}", handlers.FetchAccessPolicy).Methods(http.MethodGet)
	app.Router.HandleFunc("/access_policies", handlers.CreateOrUpdateAccessPolicy).Methods(http.MethodPost)
	app.Router.HandleFunc("/access_policies/{id}", handlers.DeleteAccessPolicy).Methods(http.MethodDelete)
}

func (app *App) initializeWebRoutes(authHandler *handlers.AuthHandler) {
	webHandler := handlers.NewWebHandler(authHandler)

	app.Router.HandleFunc("/ui/index", webHandler.IndexPage).Methods(http.MethodGet)
	app.Router.HandleFunc("/ui/user", webHandler.UserPage).Methods(http.MethodGet)
}

func (app *App) initializeUsageRoutes() {
	app.Router.HandleFunc("/usages", handlers.FetchUsages).Methods(http.MethodGet)
	app.Router.HandleFunc("/usages", handlers.CreateOrUpdateUsage).Methods(http.MethodPost)
}
