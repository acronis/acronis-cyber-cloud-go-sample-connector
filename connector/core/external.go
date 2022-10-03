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

package core

import (
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
)

// ExternalSystemClient is an interface to communicate with external system.
// Communication could be in the form of pushing change events from connector or pulling data from external system.
type ExternalSystemClient interface {
	// 1. Tenants-related changes (sync to external-system)
	// 1a. When a tenant is created or updated on Acronis cloud, connector will call CreateOrUpdateTenant
	CreateOrUpdateTenant(tenant *accclient.Tenant) (created bool, err error)

	// 1b. When a tenant is deleted on Acronis cloud, connector will call DeleteTenant and provides the tenantID
	DeleteTenant(tenantID string) error

	// 1c. For reconciliation purpose, connector needs to get list of tenantIDs
	// which are still active in external-system.
	GetActiveTenantIDs(offset, limit int) (tenantIDs []string, err error)

	// 1d. During tenant creation, we need to maintain hierarchy of tenants. In order to ensure this hierarchy,
	// connector needs to check if a particular tenantID already exists in external-system.
	CheckTenantExist(tenantID string) (tenantExists bool, err error)

	// 2. OfferingItem-related changes (sync to external-system)
	// 2a. When an offering item is enabled or updated, connector will call CreateOrUpdateOfferingItem
	CreateOrUpdateOfferingItem(item *accclient.OfferingItem) (created bool, err error)

	// 2b. When an offering item is disabled, connector will call DeleteOfferingItem
	// OfferingItemID contains the offering item name and the respective tenantID
	DeleteOfferingItem(itemID OfferingItemID) error

	// 2c. For reconciliation purpose, connector needs to get list of offering item IDs
	// which are still active in external-system.
	GetActiveOfferingItemIDs(offset, limit int) ([]OfferingItemID, error)

	// 3. User-related changes (sync to external-system)
	// 3a. When a user is created or updated, connector will call CreateOrUpdateUser
	CreateOrUpdateUser(user *accclient.User) (created bool, err error)

	// 3b. When a user is deleted, connector will cal DeleteUser and provides the userID
	DeleteUser(userID string) error

	// 3c. For reconciliation purpose, connector needs to get list of user IDs
	// which are still active in external-system.
	GetActiveUserIDs(offset, limit int) (userIDs []string, err error)

	// 4. AccessPolicy-related changes (sync to external-system)
	// 4a. When an access policy is assigned or updated for an user, connector will call CreateOrUpdateAccessPolicy
	CreateOrUpdateAccessPolicy(accessPolicy *accclient.AccessPolicy) (created bool, err error)

	// 4b. When an access policy is revoked from an user, connector will call DeleteAccessPolicy
	DeleteAccessPolicy(accessPolicyID string) error

	// 4c. For reconciliation purpose, connector needs to get list of access policy IDs
	// which are still active in external-system.
	GetActiveAccessPolicyIDs(offset, limit int) (accessPolicyIDs []string, err error)

	// 5. external-system usage-related changes (sync from external-system)
	// Connector will pull usage information from external-system to be pushed to Acronis cloud.
	GetUsages(offset, limit int) ([]accclient.Usage, error)
}

// OfferingItemID is the minimal structure that identifies an offering item uniquely on Acronis cloud
type OfferingItemID struct {
	OfferingItemName string
	TenantID         string
}
