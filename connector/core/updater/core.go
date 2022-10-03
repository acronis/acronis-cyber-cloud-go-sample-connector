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

package updater

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/core"
)

// Updater provides the functionality of syncing between Acronis Cyber Cloud Platform and external systems (ISV)
type Updater struct {
	config *Config
	sync   core.SyncLoop
	recon  core.Reconciliation
	usage  core.UsageLoop
}

type Option func(*Updater)

// NewUpdater returns an Updater initialized with the given params
func NewUpdater(config *Config, externalClient core.ExternalSystemClient, options ...Option) (*Updater, error) {
	u := &Updater{
		config: config,
	}

	for _, option := range options {
		option(u)
	}

	// Setup ACC Client
	httpDefaultClient := &http.Client{Timeout: 120 * time.Second}
	httpClient := getHTTPClient(
		config.AuthSettings.ClientID,
		config.AuthSettings.ClientSecret,
		config.APIServerSettings.BaseURL+"/api/2",
		httpDefaultClient,
	)

	accClient := accclient.NewClient(httpClient, config.APIServerSettings.BaseURL)

	tenantID, err := accClient.GetRegistrationTenantID(
		context.Background(), config.APIServerSettings.BaseURL, config.AuthSettings.ClientID)
	if err != nil {
		return nil, err
	}

	u.sync = NewSyncLoop(
		accClient,
		tenantID,
		externalClient,
		WithUpdateInterval(config.UpdateInterval),
	)

	u.recon = NewReconciliationLoop(
		accClient,
		tenantID,
		externalClient,
		WithReconciliationInterval(config.ReconciliationInterval),
	)

	u.usage = NewUsageLoop(
		accClient,
		externalClient,
		WithUsageUpdateInterval(config.UsageReportInterval),
	)

	return u, nil
}

// Start runs the update and reconciliation loops
func (u *Updater) Start() error {
	// first reconcile all objects on startup
	tenantsUpdateTimestamp := u.recon.ReconcileTenantsAndOfferingItems(true)
	usersUpdateTimestamp := u.recon.ReconcileUsersAndAccessPolicies(true)

	// run all SYNC loops
	go u.sync.UpdateTenantsAndOfferingItems(tenantsUpdateTimestamp)
	go u.sync.UpdateUsersAndAccessPolicies(usersUpdateTimestamp)
	go u.recon.ReconcileTenantsAndOfferingItems(false)
	go u.recon.ReconcileUsersAndAccessPolicies(false)
	go u.usage.UpdateUsages()

	return nil
}

// getHTTPClient returns a HTTP clent for identification with service's access token
func getHTTPClient(clientID, clientSecret, idpAddr string, httpClient *http.Client) *http.Client {
	oauth2Config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     idpAddr + "/idp/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	return oauth2Config.Client(ctx)
}
