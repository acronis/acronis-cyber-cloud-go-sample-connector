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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
)

// Example code to invoke Acronis Cloud Platform API to retrieve tenants information
const (
	ClientID     = "" // ClientID obtained from registration to Acronis Platform
	ClientSecret = "" // ClientSecret obtained from registration to Acronis Platform
	DCURL        = "" // provide Acronis Cloud DC URL here
)

// getHTTPClient returns a HTTP client for identification with service's access token
func getHTTPClient(ctx context.Context, clientID, clientSecret, idpAddr string, httpClient *http.Client) *http.Client {
	oauth2Config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     idpAddr + "/idp/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	return oauth2Config.Client(ctx)
}

func main() {
	fmt.Println("Sample Connector")

	ctx := context.Background()

	httpDefaultClient := &http.Client{Timeout: 120 * time.Second}
	httpClient := getHTTPClient(
		ctx,
		ClientID,
		ClientSecret,
		DCURL+"/api/2",
		httpDefaultClient,
	)

	client := accclient.NewClient(httpClient, DCURL)

	// Get tenant ID
	tenantID, err := client.GetRegistrationTenantID(ctx, DCURL, ClientID)
	if err != nil {
		log.Fatalf("Failed to get tenant ID: %v", err)
	}

	var limit uint = 2
	getTenantReq := &accclient.TenantGetRequest{
		SubTreeRootID: tenantID,
		LevelOfDetail: accclient.TenantLODFull,
		Limit:         &limit,
	}

	tenantResp, err := client.GetTenants(ctx, getTenantReq)
	if err != nil {
		log.Printf("Error %v", err)
		return
	}
	for {
		log.Printf("%v", tenantResp)
		log.Printf("Total number of tenant %d of this page", len(tenantResp.Items))

		tenantResp, err = client.GetTenantsNextPage(ctx, tenantResp)
		// If there is error on request
		if err != nil {
			log.Printf("Error %v", err)
			break
		}
		// If there is no more pages
		if tenantResp == nil && err == nil {
			log.Print("No more page")
			break
		}
	}
}
