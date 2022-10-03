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

// Usage represents the input params for a single Usages Update
type Usage struct {
	// ResourceID required for per-resource reporting
	ResourceID *string `json:"resource_id,omitempty"`

	// UsageType required for per-resource reporting
	UsageType *string `json:"usage_type,omitempty"`

	// TenantID required for per-tenant reporting
	TenantID *string `json:"tenant_id,omitempty"`

	// OfferingItem required for per-tenant reporting
	OfferingItem *string `json:"offering_item,omitempty"`

	// InfraID is an UUID required for per-tenant reporting for infra offering items
	InfraID *string `json:"infra_id,omitempty"`

	UsageValue int64 `json:"usage_value"`
}

// UsagesPutRequest represents the input params for the Put Usages request
type UsagesPutRequest struct {
	Items []Usage `json:"items"`
}

// UsagesResponse represents the response from PUT Usages API of a single usage
type UsagesResponse struct {
	// ResourceID required for per-resource reporting
	ResourceID *string `json:"resource_id,omitempty"`

	// UsageType required for per-resource reporting
	UsageType *string `json:"usage_type,omitempty"`

	// TenantID required for per-tenant reporting
	TenantID *string `json:"tenant_id,omitempty"`

	// OfferingItem required for per-tenant reporting
	OfferingItem *string `json:"offering_item,omitempty"`

	// InfraID is an UUID required for per-tenant reporting for infra offering items
	InfraID *string `json:"infra_id,omitempty"`

	Error *Error `json:"error,omitempty"`
}

// UsagesPutResponse represents the response from the PUT Usages API
type UsagesPutResponse struct {
	Response
	Items []UsagesResponse `json:"items"`
}

func parseUsagesPutResponse(r *http.Response) (*UsagesPutResponse, error) {
	var u UsagesPutResponse
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, err
	}

	u.StatusCode = r.StatusCode
	u.HTTPHeader = r.Header

	return &u, nil
}
