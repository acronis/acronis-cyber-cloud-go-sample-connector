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
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OfferingItem struct {
	ApplicationID string  `json:"application_id"`
	Name          string  `json:"name"`
	Edition       *string `json:"edition,omitempty"`
	UsageName     string  `json:"usage_name"`
	TenantID      string  `json:"tenant_id"`

	// Last update timestamp, if tenant has just been created - then is equal to created_at,
	// if tenant has just been deleted - then is equal to deleted_at
	UpdatedAt time.Time `json:"updated_at"`

	// Status of offering item: 1 - item turned on, 0 - off
	Status int `json:"status"`

	// Flag, if 'true' this item status can not be changed
	Locked bool `json:"locked"`

	Quota   Quota  `json:"quota"`
	Type    string `json:"type"`
	InfraID string `json:"infra_id"`

	// Measurement unit in which offering item's usages are kept (e.g.: 'bytes', 'quantity', 'seconds', 'n/a')
	MeasurementUnit string `json:"measurement_unit"`
}

type Quota struct {
	Value   *float64 `json:"value,omitempty"`
	Overage *float64 `json:"overage,omitempty"`
	Version float64  `json:"version"`
}

// OfferingItemsGetRequest represents the input params for the Get Offering Items request
type OfferingItemsGetRequest struct {
	// SubTreeRootTenantID is a filter to fetch offering items for tenants hierarchy starting from
	// (inclusive) the specified one. Sorting by tenant level is always assumed
	SubTreeRootTenantID string

	// UsageNames filters the results by the list of usage names and returns items with matching usage names
	UsageNames []string

	// Editions filters the results by the list of editions and returns items with matching editions.
	Editions []string

	// UpdatedSince is a filter to fetch users which were updated later than the specified timestamp
	UpdatedSince *time.Time

	// Limit sets the number of elements in current users page of the response.
	Limit *uint

	// After is a cursor to fetch the next users page. The cursor encodes all the filtering and sorting arguments,
	// thus client does not need to provide all them for the next page, only cursor should be provided.
	After string
}

func (o *OfferingItemsGetRequest) getQueryParam() url.Values {
	params := url.Values{}

	if o.SubTreeRootTenantID != "" {
		params.Set("subtree_root_tenant_id", o.SubTreeRootTenantID)
	}
	if len(o.UsageNames) > 0 {
		params.Set("usage_names", strings.Join(o.UsageNames, ","))
	}
	if len(o.Editions) > 0 {
		params.Set("editions", strings.Join(o.Editions, ","))
	}
	if o.UpdatedSince != nil {
		params.Set("updated_since", o.UpdatedSince.Format(time.RFC3339))
	}
	if o.Limit != nil {
		params.Set("limit", fmt.Sprintf("%d", *o.Limit))
	}
	if o.After != "" {
		params.Set("after", o.After)
	}

	return params
}

// OfferingItemsGetResponse represents the response from the Get Offering Items API
type OfferingItemsGetResponse struct {
	Response
	Pagination
	Timestamp time.Time      `json:"timestamp"`
	Items     []OfferingItem `json:"items"`
}

func parseOfferingItemsGetResponse(r *http.Response) (*OfferingItemsGetResponse, error) {
	var o OfferingItemsGetResponse
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		return nil, err
	}

	o.StatusCode = r.StatusCode
	o.HTTPHeader = r.Header

	return &o, nil
}

// OfferingItemsTenantPutRequest represents the request for the PUT Offering Items in individual tenants
type OfferingItemsTenantPutRequest struct {
	OfferingItems []*OfferingItemTenantPut `json:"offering_items"`
}

// OfferingItemTenantPut is the individual item for the PUT offering items request body.
type OfferingItemTenantPut struct {
	// Application ID of the offering item, as seen in the /applications endpoint.
	ApplicationID string `json:"application_id"`
	Name          string `json:"name"`

	// Status of offering item: 1 - item turned on, 0 - off
	Status  *int8   `json:"status,omitempty"`
	InfraID *string `json:"infra_id,omitempty"`
	Quota   *Quota  `json:"quota,omitempty"`
}
