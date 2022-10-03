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

const userEndpointBase = "/users"

// CreateOrUpdateUser creates or updates the user in external system
func (c *Client) CreateOrUpdateUser(user *models.User) (bool, error) {
	resp, err := c.doPost(userEndpointBase, *user)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	isCreated := resp.StatusCode == http.StatusCreated
	return isCreated, nil
}

// DeleteUser deletes the user by id in external system
func (c *Client) DeleteUser(id string) error {
	resp, err := c.doDelete(userEndpointBase + "/" + id)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

// GetUsers gets the list of users
func (c *Client) GetUsers(offset, limit int) ([]models.User, error) {
	apiPath := userEndpointBase + "?offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	resp, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []models.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser get user by id
func (c *Client) GetUser(id string) (*models.User, error) {
	resp, err := c.doGet(userEndpointBase + "/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	user := &models.User{}
	if err := decoder.Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}
