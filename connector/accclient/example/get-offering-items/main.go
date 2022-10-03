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
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
)

// Example code to invoke Acronis Cloud Platform API to retrieve offering items information
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

	getOfferingItemsReq := &accclient.OfferingItemsGetRequest{
		SubTreeRootTenantID: tenantID,
	}
	offeringItemsResp, err := client.GetOfferingItems(ctx, getOfferingItemsReq)
	if err != nil {
		log.Printf("Error %v", err)
		return
	}
	log.Printf("%v", offeringItemsResp)
}
