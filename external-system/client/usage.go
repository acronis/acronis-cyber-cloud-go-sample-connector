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
	"net/http"
	"strconv"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

const usagesEndpointBase = "/usages"

// GetUsages gets the list of usages
func (c *Client) GetUsages(offset, limit int) ([]models.Usage, error) {
	apiPath := usagesEndpointBase + "?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	resp, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var usages []models.Usage
	if err := json.NewDecoder(resp.Body).Decode(&usages); err != nil {
		return nil, err
	}

	return usages, nil
}

// CreateOrUpdateUsage creates or updates the usage in external system
func (c *Client) CreateOrUpdateUsage(usage *models.Usage) (bool, error) {
	resp, err := c.doPost(usagesEndpointBase, *usage)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	isCreated := resp.StatusCode == http.StatusCreated
	return isCreated, nil
}
