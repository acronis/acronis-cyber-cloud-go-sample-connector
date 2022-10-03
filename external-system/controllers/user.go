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

package controllers

import (
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// GetUsers fetches list of users
func GetUsers(offset, limit int) ([]models.User, error) {
	return models.GetUsers(offset, limit)
}

// GetUser fetches user by id
func GetUser(id string) (models.User, error) {
	return models.GetUser(id)
}

// CreateOrUpdateUser Updates existing user or creates the user if it doesnt exist
func CreateOrUpdateUser(user *models.User) (bool, error) {
	return models.CreateOrUpdateUser(user)
}

// DeleteUser deletes a single user
func DeleteUser(id string) error {
	return models.DeleteUser(id)
}
