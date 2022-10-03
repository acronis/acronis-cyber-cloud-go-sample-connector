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

// The `controllers` package provides actions to access data layer
package controllers

import (
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// Get the list of access policies from data access layer
func GetAccessPolicies(offset, limit int) ([]models.AccessPolicy, error) {
	return models.GetAccessPolicies(offset, limit)
}

/* CreateOrUpdateAccessPolicy function creates new access policy if it doesn't exist
   else update the existing access policy */
func CreateOrUpdateAccessPolicy(accessPolicy *models.AccessPolicy) (bool, error) {
	return models.CreateOrUpdateAccessPolicy(accessPolicy)
}

// Get the access policy from data access layer
func GetAccessPolicy(id string) (*models.AccessPolicy, error) {
	return models.GetAccessPolicy(id)
}

// Delete access policy in external system db
func DeleteAccessPolicy(id string) error {
	return models.DeleteAccessPolicy(id)
}
