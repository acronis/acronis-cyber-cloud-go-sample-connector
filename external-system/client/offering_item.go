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

// CreateOrUpdateOfferingItem creates or updates the offering items in external system
func (c *Client) CreateOrUpdateOfferingItem(oi *models.OfferingItem) (bool, error) {
	apiPath := "/tenants/" + oi.TenantID + "/offering_items"
	resp, err := c.doPost(apiPath, oi)
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

// DeleteOfferingItem deletes the offering item
func (c *Client) DeleteOfferingItem(tenantID, name string) error {
	apiPath := "/tenants/" + tenantID + "/offering_items/" + name
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

// GetOfferingItems gets the list of offering items
func (c *Client) GetOfferingItems(offset, limit int) ([]models.OfferingItem, error) {
	apiPath := "/offering_items"
	apiPath = apiPath + "?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	resp, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error status code %v returned by server", resp.StatusCode)
	}

	var items []models.OfferingItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetOfferingItem get offeringItem by name
func (c *Client) GetOfferingItem(tenantID, name string) (*models.OfferingItem, error) {
	apiPath := "/tenants/" + tenantID + "/offering_items/" + name
	resp, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error status code %v returned by server", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	oi := &models.OfferingItem{}
	if err := decoder.Decode(oi); err != nil {
		return nil, err
	}

	return oi, nil
}
