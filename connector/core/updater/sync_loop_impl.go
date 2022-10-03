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

package updater

import (
	"context"
	"time"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/core"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/logs"
)

// SyncLoopImpl is a sample implementation of connector to sync items from
// Acronis Cyber Cloud Platform to external-system (ISV)
type SyncLoopImpl struct {
	accClient *accclient.Client
	tenantID  string
	extClient core.ExternalSystemClient

	// optional to be set during initialization
	updateInterval uint // in seconds

	// last response time from Acronis cloud for tenants and offering items update loop
	tenantsLoopUpdatedSince *time.Time

	// last response time from Acronis cloud for users and access policies update loop
	usersLoopUpdatedSince *time.Time
}

// NewSyncLoop initializes SyncLoopImpl as an implementation of core.SyncLoop
func NewSyncLoop(
	accClient *accclient.Client,
	tenantID string,
	extClient core.ExternalSystemClient,
	options ...func(*SyncLoopImpl)) core.SyncLoop {
	syncLoop := &SyncLoopImpl{
		accClient:      accClient,
		tenantID:       tenantID,
		extClient:      extClient,
		updateInterval: 5, // default
	}

	for _, option := range options {
		option(syncLoop)
	}
	return syncLoop
}

// WithUpdateInterval is an optional init function to set update interval
func WithUpdateInterval(updateInterval uint) func(*SyncLoopImpl) {
	return func(loop *SyncLoopImpl) {
		loop.updateInterval = updateInterval
	}
}

// UpdateTenantsAndOfferingItems syncs Tenants and Offering Items changes from
// Acronis Cyber Cloud Platform to external-system
// 1. Pulls tenants and offering items changes with updated_since filter
// 2. For each tenant:
//	  a. If tenantID exists and deletedAt is empty, perform create or update (upsert)
//    b. Create, Update or Delete offering items
//    c. If tenantID exists and deletedAt is non-empty, perform delete
//    d. If tenantID doesn't exist, perform delete via TenantID in OfferingItems
func (loop *SyncLoopImpl) UpdateTenantsAndOfferingItems(firstUpdatedSince time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.ContextID, "tenants_loop")
	logger := logs.GetDefaultLogger(ctx)

	if !firstUpdatedSince.IsZero() {
		loop.tenantsLoopUpdatedSince = &firstUpdatedSince
	}

	limit := uint(100)
	withContacts := true
	withOfferingItems := true
	for ; ; time.Sleep(time.Second * time.Duration(loop.updateInterval)) {
		tenantsRequest := &accclient.TenantGetRequest{
			SubTreeRootID:     loop.tenantID,
			Limit:             &limit,
			WithContacts:      &withContacts,
			WithOfferingItems: &withOfferingItems,
			AllowDeleted:      true,
			UpdatedSince:      loop.tenantsLoopUpdatedSince,
		}

		tenantsResp, err := loop.accClient.GetTenants(ctx, tenantsRequest)
		if err != nil {
			logger.Warnf("Failed to get tenants: %v", err)
			continue
		}

		loop.tenantsLoopUpdatedSince = &tenantsResp.Timestamp.Time

		syncedTenantsCount := len(tenantsResp.Items)
		syncedOfferingItemsCount := countOfferingItemsInTenants(tenantsResp.Items)
		loop.processTenantsAndOfferingItemsChanges(ctx, tenantsResp.Items)

		for tenantsResp != nil && tenantsResp.After() != "" {
			tenantsResp, err = loop.accClient.GetTenantsNextPage(ctx, tenantsResp)
			if err != nil {
				logger.Warnf("Failed to get tenants next page: %v", err)
				break
			}

			if tenantsResp != nil {
				syncedTenantsCount += len(tenantsResp.Items)
				syncedOfferingItemsCount += countOfferingItemsInTenants(tenantsResp.Items)
				loop.processTenantsAndOfferingItemsChanges(ctx, tenantsResp.Items)
			}
		}

		if syncedTenantsCount > 0 {
			logger.Infof("Synced %v tenants and %v offering items", syncedTenantsCount, syncedOfferingItemsCount)
		} else {
			logger.Debug("Update loop succeed, no tenants changes reported")
		}
	}
}

// UpdateUsersAndAccessPolicies syncs Users and access policies changes from
// Acronis Cyber Cloud Platform to external-system
// 1. Pulls users and access policies changes with updated_since filter
// 2. For each user:
//    a. If userID exists and deletedAt is empty, perform create or update (upsert)
//    b. Create, Update or Delete access policies
//    c. If userID exists and deletedAt is non-empty, perform delete
//    d. If userID doesn't exist, perform delete via TrusteeID in access policy
func (loop *SyncLoopImpl) UpdateUsersAndAccessPolicies(firstUpdatedSince time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.ContextID, "users_loop")
	logger := logs.GetDefaultLogger(ctx)

	if !firstUpdatedSince.IsZero() {
		loop.usersLoopUpdatedSince = &firstUpdatedSince
	}

	limit := uint(100)
	withAccessPolicies := true
	for ; ; time.Sleep(time.Second * time.Duration(loop.updateInterval)) {
		usersRequest := &accclient.UserGetRequest{
			SubTreeRootTenantID: loop.tenantID,
			Limit:               &limit,
			WithAccessPolicies:  &withAccessPolicies,
			AllowDeleted:        true,
			UpdatedSince:        loop.usersLoopUpdatedSince,
		}

		usersResp, err := loop.accClient.GetUsers(ctx, usersRequest)
		if err != nil {
			logger.Warnf("Failed to get users: %v", err)
			continue
		}

		loop.usersLoopUpdatedSince = &usersResp.Timestamp

		syncedUsersCount := len(usersResp.Items)
		syncedAccessPoliciesCount := countAccessPoliciesInUsers(usersResp.Items)
		loop.processUsersAndAccessPoliciesChanges(ctx, usersResp.Items)

		for usersResp != nil && usersResp.After() != "" {
			usersResp, err = loop.accClient.GetUsersNextPage(ctx, usersResp)
			if err != nil {
				logger.Warnf("Failed to get users next page: %v", err)
				break
			}

			if usersResp != nil {
				syncedUsersCount += len(usersResp.Items)
				syncedAccessPoliciesCount += countAccessPoliciesInUsers(usersResp.Items)
				loop.processUsersAndAccessPoliciesChanges(ctx, usersResp.Items)
			}
		}

		if syncedUsersCount > 0 {
			logger.Infof("Synced %v users and %v access policies", syncedUsersCount, syncedAccessPoliciesCount)
		} else {
			logger.Debug("Update loop succeed, no users changes reported")
		}
	}
}

// ===================
// helper functions
// ===================

// processTenantsAndOfferingItemsChanges processes changes reported by composite API of tenants and offering items.
func (loop *SyncLoopImpl) processTenantsAndOfferingItemsChanges(ctx context.Context, items []accclient.Tenant) {
	logger := logs.GetDefaultLogger(ctx)

	for i := range items {
		deleteTenantID := ""
		if items[i].ID != "" {
			if items[i].DeletedAt.IsZero() {
				if err := createOrUpdateTenant(ctx, loop.extClient, loop.accClient, &items[i]); err != nil {
					// error is treated as non-fatal, skip and continue to next tenant
					logger.Warnf("Failed to update tenant %v: %s", items[i].ID, err)
				}
			} else {
				deleteTenantID = items[i].ID
			}
		} else if len(items[i].OfferingItems) > 0 {
			deleteTenantID = items[i].OfferingItems[0].TenantID
		}

		loop.processOfferingItemsChanges(ctx, items[i].OfferingItems)

		// perform tenant deletion after processing offering items
		if deleteTenantID != "" {
			if err := loop.extClient.DeleteTenant(deleteTenantID); err != nil {
				logger.Warnf("Failed to push tenant deletion to external system: %v", err)
			}
		}
	}
}

// processOfferingItemsChanges pushes offering items change events to external system
func (loop *SyncLoopImpl) processOfferingItemsChanges(ctx context.Context, items []accclient.OfferingItem) {
	logger := logs.GetDefaultLogger(ctx)
	for i := range items {
		if items[i].Status == 0 {
			if err := loop.extClient.DeleteOfferingItem(core.OfferingItemID{
				OfferingItemName: items[i].Name,
				TenantID:         items[i].TenantID,
			}); err != nil {
				logger.Warnf("Failed to delete offering item %v for tenant %v: %v", items[i].Name, items[i].TenantID, err)
			}
		} else {
			if oiCreated, err := loop.extClient.CreateOrUpdateOfferingItem(&items[i]); err != nil {
				logger.Warnf("Failed to upsert offering item %v for tenant %v into external-system: %v",
					items[i].Name, items[i].TenantID, err)
			} else {
				logger.Debugf("Offering item %v for tenant %v successfully updated (is new offering item: %v)",
					items[i].Name, items[i].TenantID, oiCreated)
			}
		}
	}
}

// processUsersAndAccessPoliciesChanges processes changes reported by composite API of users and access policies
func (loop *SyncLoopImpl) processUsersAndAccessPoliciesChanges(ctx context.Context, items []accclient.User) {
	logger := logs.GetDefaultLogger(ctx)

	for i := range items {
		deleteUserID := ""
		// ID field exists if user has active access policies
		if items[i].ID != "" {
			if items[i].DeletedAt.IsZero() {
				if err := createOrUpdateUser(ctx, loop.extClient, loop.accClient, &items[i]); err != nil {
					// error is treated as non-fatal, skip and continue to next user
					logger.Warnf("Failed to update user %v: %s", items[i].ID, err)
				}
			} else {
				deleteUserID = items[i].ID
			}
		} else if len(items[i].AccessPolicies) > 0 {
			deleteUserID = items[i].AccessPolicies[0].TrusteeID
		}

		loop.processAccessPoliciesChanges(ctx, items[i].AccessPolicies)

		// perform user deletion after processing access policies
		if deleteUserID != "" {
			if err := loop.extClient.DeleteUser(deleteUserID); err != nil {
				logger.Warnf("Failed to push user deletion to external system: %v", err)
			}
		}
	}
}

// processAccessPoliciesChanges pushes access policies change events to external system
func (loop *SyncLoopImpl) processAccessPoliciesChanges(ctx context.Context, items []accclient.AccessPolicy) {
	logger := logs.GetDefaultLogger(ctx)
	for i := range items {
		if items[i].DeletedAt != nil {
			if err := loop.extClient.DeleteAccessPolicy(items[i].ID); err != nil {
				logger.Warnf("Failed to delete access policy: %v", err)
			}
		} else {
			if apCreated, err := loop.extClient.CreateOrUpdateAccessPolicy(&items[i]); err != nil {
				logger.Warnf("Failed to upsert access policy %v with ID %v for user %v into external-system: %v",
					items[i].RoleID, items[i].ID, items[i].TrusteeID, err)
			} else {
				logger.Debugf("Access policy %v with ID %v for user %v successfully updated (is new access policy: %v)",
					items[i].RoleID, items[i].ID, items[i].TrusteeID, apCreated)
			}
		}
	}
}

func countOfferingItemsInTenants(items []accclient.Tenant) (counter uint) {
	for i := range items {
		counter += uint(len(items[i].OfferingItems))
	}
	return
}

func countAccessPoliciesInUsers(items []accclient.User) (counter uint) {
	for i := range items {
		counter += uint(len(items[i].AccessPolicies))
	}
	return
}
