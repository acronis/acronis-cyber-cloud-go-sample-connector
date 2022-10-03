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

// CreateOrUpdateAccessPolicy creates or updates the access policy in external system
func (c *Client) CreateOrUpdateAccessPolicy(accessPolicy *models.AccessPolicy) (bool, error) {
	const apiPath = "/access_policies"
	resp, err := c.doPost(apiPath, accessPolicy)
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

// DeleteAccessPolicy deletes the access policy
func (c *Client) DeleteAccessPolicy(id string) error {
	apiPath := "/access_policies/" + id
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

// GetAccessPolicies gets the list of access policies
func (c *Client) GetAccessPolicies(offset, limit int) ([]models.AccessPolicy, error) {
	apiPath := "/access_policies"
	apiPath = apiPath + "?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	resp, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error status code %v returned by server", resp.StatusCode)
	}

	var items []models.AccessPolicy
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetAccessPolicy get access policy by id
func (c *Client) GetAccessPolicy(id string) (*models.AccessPolicy, error) {
	apiPath := "/access_policies/" + id
	resp, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error status code %v returned by server", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	accessPolicy := &models.AccessPolicy{}
	if err := decoder.Decode(accessPolicy); err != nil {
		return nil, err
	}

	return accessPolicy, nil
}

// CheckAccessPolicyExist checks if a access policy with the specified id exists in the external-system.
func (c *Client) CheckAccessPolicyExist(id string) (bool, error) {
	apiPath := c.APIURL + "/acess_policies/" + id
	resp, err := c.doGet(apiPath)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}
