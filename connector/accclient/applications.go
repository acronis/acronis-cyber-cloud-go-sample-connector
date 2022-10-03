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

package accclient

import (
	"encoding/json"
	"net/http"
)

// Application is the individual item for the GET application response body.
type Application struct {
	ID string `json:"id"`

	// Human-readable name that will be displayed to the users
	Name string `json:"name"`

	// See RAML for updated list of supported types
	Type string `json:"type"`

	// Available usages for this application
	Usages []string `json:"usages"`

	// Base URL for application's HTTP API
	APIBaseURL string `json:"api_base_url"`
}

// ApplicationsGetResponse represents the response from the GET applications endpoint
type ApplicationsGetResponse struct {
	Response
	Items []*Application `json:"items"`
}

func parseApplicationsGetResponse(r *http.Response) (*ApplicationsGetResponse, error) {
	var t ApplicationsGetResponse
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}

	t.StatusCode = r.StatusCode
	t.HTTPHeader = r.Header

	return &t, nil
}
