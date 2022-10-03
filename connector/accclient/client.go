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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const clientRegistrationInfoEndpoint = "/api/2/clients/%s"

// Client is a Acronis Cyber Cloud Platform API client. Create one by calling NewClient
type Client struct {
	APIURL     string
	HTTPClient *http.Client
}

// NewClient creates a new Acronis Cyber Cloud Platform API client
func NewClient(httpClient *http.Client, url string) *Client {
	return &Client{
		APIURL:     url + "/api/2",
		HTTPClient: httpClient,
	}
}

// DoGet sends a HTTP GET request with the given params
func (c *Client) DoGet(ctx context.Context, url string) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, url, nil, nil)
}

// DoPost sends a HTTP POST request with the given params
func (c *Client) DoPost(ctx context.Context, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Do(ctx, http.MethodPost, url, body, headers)
}

// DoPut sends a HTTP PUT request with the given params
func (c *Client) DoPut(ctx context.Context, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Do(ctx, http.MethodPut, url, body, headers)
}

// DoDelete sends a HTTP DELETE request with the given params
func (c *Client) DoDelete(ctx context.Context, url string) (*http.Response, error) {
	return c.Do(ctx, http.MethodDelete, url, nil, nil)
}

// Do sends a HTTP request with the given params
func (c *Client) Do(ctx context.Context, method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reader io.Reader

	if body != nil {
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to create http request body. %w", err)
		}

		reader = bytes.NewBuffer(reqBody)
	} else {
		reader = nil
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request. %w", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error connecting to api server. %w", err)
	}

	if resp.StatusCode >= 400 {
		defer CloseBody(resp)
		return resp, makeError(resp)
	}

	return resp, nil
}

func CloseBody(r *http.Response) {
	if r != nil {
		// Ensure that we read until the response is complete AND call Close()
		// Do not have to check for error as the execution can continue
		_, _ = io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}
}

// GetRegistrationTenantID will retrieve the tenantID used during the application registration.
func (c *Client) GetRegistrationTenantID(ctx context.Context, baseURL, clientID string) (string, error) {
	apiPath := baseURL + fmt.Sprintf(clientRegistrationInfoEndpoint, clientID)

	resp, err := c.DoGet(ctx, apiPath)
	if err != nil {
		return "", fmt.Errorf("error in http request getRegistrationTenantID. %w", err)
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	var registrationInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&registrationInfo); err != nil {
		return "", fmt.Errorf("error parsing response getRegistrationTenantID. %w", err)
	}

	if tenantID, ok := registrationInfo["tenant_id"].(string); ok {
		return tenantID, nil
	}

	return "", errors.New("no tenantID in registration information found")
}

// GetTenants gets the list of tenants details filter by the request params
func (c *Client) GetTenants(ctx context.Context, getReq *TenantGetRequest) (*TenantGetResponse, error) {
	apiPath := c.APIURL + "/tenants"
	apiPath += "?" + getReq.getQueryParam().Encode()

	resp, err := c.DoGet(ctx, apiPath)
	if err != nil {
		return nil, fmt.Errorf("error in http request GetTenants. %w", err)
	}
	defer CloseBody(resp)

	tenantResp, err := parseTenantGetResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response GetTenants. %w", err)
	}

	if tenantResp.StatusCode != http.StatusOK {
		return tenantResp, fmt.Errorf("invalid status code %d from GetTenants", tenantResp.StatusCode)
	}

	return tenantResp, nil
}

// GetTenantsNextPage retrieves the next page of this response and returns a new TenantGetResponse.
// If both error and *TenantGetResponse are nil, there is no more page available
func (c *Client) GetTenantsNextPage(ctx context.Context, page Page) (*TenantGetResponse, error) {
	if page.After() != "" {
		req := &TenantGetRequest{
			After: page.After(),
		}
		resp, err := c.GetTenants(ctx, req)
		if err != nil {
			return resp, fmt.Errorf("error fetching GetTenantsNextPage with cursor %s. %w", page.After(), err)
		}
		return resp, nil
	}

	return nil, nil
}

// GetOfferingItems gets the list of offering items details filter by the request params
func (c *Client) GetOfferingItems(ctx context.Context, getReq *OfferingItemsGetRequest) (*OfferingItemsGetResponse, error) {
	apiPath := c.APIURL + "/tenants/offering_items"
	apiPath += "?" + getReq.getQueryParam().Encode()

	resp, err := c.DoGet(ctx, apiPath)
	if err != nil {
		return nil, fmt.Errorf("error in http request GetOfferingItems. %w", err)
	}
	defer CloseBody(resp)

	offeringItemsResp, err := parseOfferingItemsGetResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response GetOfferingItems. %w", err)
	}

	if offeringItemsResp.StatusCode != http.StatusOK {
		return offeringItemsResp, fmt.Errorf("invalid status code %d from GetOfferingItems", offeringItemsResp.StatusCode)
	}

	return offeringItemsResp, nil
}

// GetOfferingItemsNextPage retrieves the next page of this response and returns a new OfferingItemsGetResponse.
// If both error and *OfferingItemsGetResponse are nil, there is no more page available
func (c *Client) GetOfferingItemsNextPage(ctx context.Context, page Page) (*OfferingItemsGetResponse, error) {
	if page.After() != "" {
		req := &OfferingItemsGetRequest{
			After: page.After(),
		}
		resp, err := c.GetOfferingItems(ctx, req)
		if err != nil {
			return resp, fmt.Errorf("error fetching OfferingItemsGetResponse next page with cursor %s. %w", page.After(), err)
		}
		return resp, nil
	}

	return nil, nil
}

// UpdateTenantOfferingItems updates the offering items available to the tenant
func (c *Client) UpdateTenantOfferingItems(ctx context.Context, tenantID string, putReq *OfferingItemsTenantPutRequest) error {
	apiPath := fmt.Sprintf("%v/tenants/%v/offering_items", c.APIURL, tenantID)
	resp, err := c.DoPut(ctx, apiPath, putReq, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return fmt.Errorf("error in http request UpdateTenantOfferingItems. %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	return nil
}

// GetApplications gets a list of all applications available on the server
func (c *Client) GetApplications(ctx context.Context) (*ApplicationsGetResponse, error) {
	apiPath := fmt.Sprintf("%v/applications", c.APIURL)
	resp, err := c.DoGet(ctx, apiPath)
	if err != nil {
		return nil, fmt.Errorf("error in http request GetApplications. %w", err)
	}
	defer CloseBody(resp)

	applications, err := parseApplicationsGetResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response GetApplications. %w", err)
	}
	return applications, nil
}

// GetUsers gets the list of users filter by the request params
func (c *Client) GetUsers(ctx context.Context, getReq *UserGetRequest) (*UserGetResponse, error) {
	apiPath := c.APIURL + "/users"
	apiPath += "?" + getReq.getQueryParam().Encode()

	resp, err := c.DoGet(ctx, apiPath)
	if err != nil {
		return nil, fmt.Errorf("error in http request GetUsers. %w", err)
	}
	defer CloseBody(resp)

	usersResp, err := parseUserGetResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response GetUsers. %w", err)
	}

	if usersResp.StatusCode != http.StatusOK {
		return usersResp, fmt.Errorf("invalid status code %d from GetUsers", usersResp.StatusCode)
	}

	return usersResp, nil
}

// GetUsersNextPage retrieves the next page of this response and returns a new UserGetResponse.
// If both error and *UserGetResponse are nil, there is no more page available
func (c *Client) GetUsersNextPage(ctx context.Context, page Page) (*UserGetResponse, error) {
	if page.After() != "" {
		req := &UserGetRequest{
			After: page.After(),
		}
		resp, err := c.GetUsers(ctx, req)
		if err != nil {
			return resp, fmt.Errorf("error fetching UserGetResponse next page with cursor %s. %w", page.After(), err)
		}
		return resp, nil
	}

	return nil, nil
}

// CreateUser creates a user in Acronis cloud with the given object
func (c *Client) CreateUser(ctx context.Context, user *UserPost) (*UserPostResponse, error) {
	apiPath := c.APIURL + "/users"

	resp, err := c.DoPost(ctx, apiPath, user, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, fmt.Errorf("error in http request CreateUser. %w", err)
	}
	defer CloseBody(resp)

	userResp, err := parseUserPostRespose(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response PostUser. %w", err)
	}

	if userResp.StatusCode != http.StatusOK {
		return userResp, fmt.Errorf("invalid status code %d from PostUser", userResp.StatusCode)
	}

	return userResp, nil
}

// DeleteUser deletes a user in Acronis cloud with the given userID
// First it will retrieve the latest version for the specified user before performing the deletion.
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	userRequest := UserGetRequest{UUIDs: []string{userID}}
	userResp, err := c.GetUsers(ctx, &userRequest)
	if err != nil {
		return fmt.Errorf("error in GetUsers. %w", err)
	}

	version := userResp.Items[0].Version
	apiPath := fmt.Sprintf("%v/users/%v?version=%v", c.APIURL, userID, version)

	resp, err := c.DoDelete(ctx, apiPath)
	if err != nil {
		return fmt.Errorf("error in http request DeleteUser. %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid status code %d from DeleteUser", resp.StatusCode)
	}

	return nil
}

// UpdateUsages updates the list of usages listed by the params
func (c *Client) UpdateUsages(ctx context.Context, usages *UsagesPutRequest) (*UsagesPutResponse, error) {
	apiPath := c.APIURL + "/tenants/usages"

	resp, err := c.DoPut(ctx, apiPath, usages, nil)
	if err != nil {
		return nil, fmt.Errorf("error in http request UpdateUsages. %w", err)
	}
	defer CloseBody(resp)

	usagesResp, err := parseUsagesPutResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error parsing response UpdateUsages. %w", err)
	}

	if usagesResp.StatusCode != http.StatusOK {
		return usagesResp, fmt.Errorf("invalid status code %d from UpdateUsages", usagesResp.StatusCode)
	}

	return usagesResp, nil
}

// GetTenant using the UUID tenantID, returning a Tenant object if successful.
func (c *Client) GetTenant(ctx context.Context, tenantID string) (*SingleTenantResponse, error) {
	resp, err := c.DoGet(ctx, c.APIURL+"/tenants/"+tenantID)
	if err != nil {
		return nil, fmt.Errorf("error in http request GetTenant. %w", err)
	}
	defer CloseBody(resp)
	tenant, err := parseTenantResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response GetTenant. %w", err)
	}

	return tenant, nil
}

// CreateTenant with the provided input, returning a Tenant object if successful
func (c *Client) CreateTenant(ctx context.Context, postReq *TenantPostRequest) (*SingleTenantResponse, error) {
	resp, err := c.DoPost(ctx, c.APIURL+"/tenants", postReq, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, fmt.Errorf("error in http request CreateTenant. %w", err)
	}
	defer CloseBody(resp)

	tenant, err := parseTenantResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response CreateTenant. %w", err)
	}
	return tenant, nil
}

// DeleteTenant using the UUID tenantID. The tenant must first be disabled using UpdateTenant for this to succeed.
// As the version number of the tenant is necessary to perform deletion,
// 	the function will first get the current version of the tenant.
func (c *Client) DeleteTenant(ctx context.Context, tenantID string) error {
	tenantObj, err := c.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("error getting tenant version in DeleteTenant. %w", err)
	}
	deletePath := fmt.Sprintf("%v/tenants/%v/?version=%v", c.APIURL, tenantID, tenantObj.Version)
	resp, err := c.DoDelete(ctx, deletePath)
	if err != nil {
		return fmt.Errorf("error in http request DeleteTenant. %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	return nil
}

// UpdateTenant using the UUID tenantID and newDetails.
func (c *Client) UpdateTenant(ctx context.Context, tenantID string, newDetails *TenantPutRequest) error {
	resp, err := c.DoPut(ctx, c.APIURL+"/tenants/"+tenantID, newDetails, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return fmt.Errorf("error in http request UpdateTenant. %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	return nil
}

// UpdateAccessPolicy updates access policy for a particular user
func (c *Client) UpdateAccessPolicy(
	ctx context.Context, userID string, accessPolicies *AccessPolicyList) (*UpdateAccessPolicyResponse, error) {
	updateEndpoint := fmt.Sprintf("/users/%v/access_policies", userID)
	resp, err := c.DoPut(ctx, c.APIURL+updateEndpoint, accessPolicies, map[string]string{"Content-Type": "application/json"})
	if err != nil {
		return nil, fmt.Errorf("error in http request UpdateAccessPolicy. %w", err)
	}
	defer CloseBody(resp)

	accessPolicyList, err := parseUpdateAccessPolicyResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response UpdateAccessPolicy. %w", err)
	}
	return accessPolicyList, nil
}
