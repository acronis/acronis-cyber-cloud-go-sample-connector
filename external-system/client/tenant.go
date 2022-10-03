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

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// CreateOrUpdateTenant creates or updates the tenant in external system
func (c *Client) CreateOrUpdateTenant(tenant *models.Tenant) (bool, error) {
	const apiPath = "/tenants"
	resp, err := c.doPost(apiPath, tenant)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return false, fmt.Errorf("error status code %v returned by server", resp.StatusCode)
	}

	isCreated := resp.StatusCode == http.StatusCreated
	return isCreated, nil
}

// DeleteTenant deletes the tenant
func (c *Client) DeleteTenant(id string) error {
	apiPath := "/tenants/" + id
	resp, err := c.doDelete(apiPath)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("error status code %v returned by server", resp.StatusCode)
	}

	return nil
}

// GetTenants gets the list of tenants
func (c *Client) GetTenants(offset, limit int) ([]models.Tenant, error) {
	apiPath := "/tenants"
	apiPath = apiPath + "?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	resp, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error status code %v returned by server", resp.StatusCode)
	}

	var items []models.Tenant
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetTenant get tenant by id
func (c *Client) GetTenant(id string) (*models.Tenant, error) {
	apiPath := "/tenants/" + id
	resp, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error status code %v returned by server", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	tenant := &models.Tenant{}
	if err := decoder.Decode(tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// CheckTenantExist checks if a tenant with the specified id exists in the external-system.
func (c *Client) CheckTenantExist(id string) (bool, error) {
	apiPath := "/tenants/" + id
	resp, err := c.doGet(apiPath)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}
