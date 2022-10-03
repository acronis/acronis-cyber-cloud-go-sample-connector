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

// +build e2e

package tests

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Please provide the values according to description in readme
const (
	clientID          = ""
	clientSecret      = ""
	dcURL             = ""
	applicationType   = ""
	roleName          = ""
	externalSystemURL = "http://externalsystem:8080" // update this once new implementation of external-system is provided
	updateInterval    = 10                           // update this according to update interval in connector config if it's modified there
)

var (
	offeringItemNames = []string{}
)

// ======
// STOP config modification beyond this line
// ======

func getTestHTTPClient(ctx context.Context, clientID, clientSecret, idpAddr string, httpClient *http.Client) *http.Client {
	oauth2Config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     idpAddr + "/api/2/idp/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	return oauth2Config.Client(ctx)
}

func getTimestampedName(name string) string {
	t := time.Now()
	return fmt.Sprintf("%s-%d-%02d-%02dT%02d:%02d:%02d\n",
		name, t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
