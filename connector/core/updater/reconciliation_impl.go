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
	"fmt"
	"time"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/core"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/logs"
)

const accPageSize = 100            // number of items requested per request to Acronis Cyber Cloud
const externalSystemPageSize = 100 // number of items requested per request to external system

// ReconciliationLoop is a sample implementation of connector to reconcile items between
// Acronis Cyber Cloud Platform and external-system
type ReconciliationLoop struct {
	accClient *accclient.Client
	tenantID  string
	extClient core.ExternalSystemClient
	ctx       context.Context

	// optional to be set during initialization
	reconciliationInterval uint // in seconds
}

// NewReconciliationLoop initializes ReconciliationLoop as an implementation of core.Reconciliation
func NewReconciliationLoop(
	accClient *accclient.Client,
	tenantID string,
	extClient core.ExternalSystemClient,
	options ...func(*ReconciliationLoop)) core.Reconciliation {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.ContextID, "reconciliation_loop")
	loop := &ReconciliationLoop{
		accClient:              accClient,
		tenantID:               tenantID,
		extClient:              extClient,
		ctx:                    ctx,
		reconciliationInterval: 3600, // default
	}

	for _, option := range options {
		option(loop)
	}

	return loop
}

// WithReconciliationInterval is an optional init function to set reconciliation interval
func WithReconciliationInterval(interval uint) func(*ReconciliationLoop) {
	return func(loop *ReconciliationLoop) {
		loop.reconciliationInterval = interval
	}
}

// ReconcileTenantsAndOfferingItems will sync all tenants and offering items
// between Acronis Cyber Cloud and external system periodically
// 1. Get tenants from ACC and embed offering_items into each tenant
// 2. Get tenant IDs which currently exist on external system
// 3. Remove tenant from external system if it doesn't exist in ACC anymore
// 4. For each tenant from ACC, push into external system to be created/updated (upsert operation)
// 5. Get Offering items from external system
// 6. Remove offering item from external system if it is not active in ACC anymore
// 7. For each active offering items from ACC, push into external system to be created/updated (upsert operation)
// If onStartup is set to true, it will only run the logic above once to make sure all tenants
// and offering items are in sync upon startup. It will also return timestamp that could be used as
// updated_since filter for the subsequent update loop
// If onStartup is set to false, the logic will be run periodically every reconciliationInterval in config file
func (loop *ReconciliationLoop) ReconcileTenantsAndOfferingItems(onStartup bool) time.Time {
	if onStartup {
		return loop.reconcileTenantsAndOfferingItems()
	}

	for {
		// wait for next cycle of reconciliation if it's not the first "sync" on startup
		time.Sleep(time.Second * time.Duration(loop.reconciliationInterval))
		loop.reconcileTenantsAndOfferingItems()
	}
}

func (loop *ReconciliationLoop) reconcileTenantsAndOfferingItems() time.Time {
	logger := logs.GetDefaultLogger(loop.ctx)

	// 1. Get tenants from ACC
	var accTenants map[string]*accclient.Tenant
	var nextUpdateTimestamp time.Time
	err := retryHelper(loop.ctx,
		func() error {
			var getRequestError error
			accTenants, nextUpdateTimestamp, getRequestError = loop.getACCTenantsAndOfferingItemsForReconciliation()
			return getRequestError
		})
	if err != nil {
		logger.Warnf("Failed to get ACC tenants: %v", err)
		return nextUpdateTimestamp // retry in next loop
	}

	// 2. Get tenants from External System
	var externalTenantIDs map[string]struct{}
	err = retryHelper(loop.ctx,
		func() error {
			var getRequestError error
			externalTenantIDs, getRequestError = loop.getExternalSystemTenantIDs()
			return getRequestError
		})
	if err != nil {
		logger.Warnf("Failed to get external tenants: %v", err)
		return nextUpdateTimestamp // retry in next loop
	}

	// 3. remove non-existing tenants
	for externalTenantID := range externalTenantIDs {
		if _, ok := accTenants[externalTenantID]; !ok {
			logger.Infof("Removing tenant %v", externalTenantID)
			if deleteErr := loop.extClient.DeleteTenant(externalTenantID); deleteErr != nil {
				logger.Warnf("Failed to delete tenant %v: %v", externalTenantID, deleteErr)
			}
		}
	}

	// 4. create or update tenant
	for accTenantID, accTenant := range accTenants {
		logger.Infof("Updating tenant %v", accTenantID)
		if upsertErr := createOrUpdateTenant(loop.ctx, loop.extClient, loop.accClient, accTenant); upsertErr != nil {
			logger.Warnf("Failed to update tenant %v: %v", accTenantID, upsertErr)
		}
	}

	// 5. get external offering items
	var externalOIs []core.OfferingItemID
	err = retryHelper(loop.ctx,
		func() error {
			var getRequestError error
			externalOIs, getRequestError = loop.getExternalSystemOfferingItems()
			return getRequestError
		})
	if err != nil {
		logger.Warnf("Failed to get external offering items: %v", err)
		return nextUpdateTimestamp // retry in next loop
	}

	// 6. remove non existing offering items
	loop.deleteInactiveOfferingItems(accTenants, externalOIs)

	// 7. create or update OI
	loop.updateOfferingItemsOnExternalSystem(accTenants)

	return nextUpdateTimestamp
}

// ReconcileUsersAndAccessPolicies will sync all users and access policies
// between Acronis Cyber Cloud and external system periodically
// 1. Get users from ACC and embed access policies into each user
// 2. Get user IDs which currently exist on external system
// 3. Remove user from external system if it doesn't exist in ACC anymore
// 4. For each user from ACC, push into external system to be created/updated (upsert operation)
// 5. Get active access policies from external system
// 6. Remove access policy from external system if it is not active in ACC anymore
// 7. For each access policy from ACC, push into external system to be created/updated (upsert operation)
// If onStartup is set to true, it will only run the logic above once to make sure all tenants
// and offering items are in sync upon startup. It will also return timestamp that could be used as
// updated_since filter for the subsequent update loop
// If onStartup is set to false, the logic will be run periodically every reconciliationInterval in config file
func (loop *ReconciliationLoop) ReconcileUsersAndAccessPolicies(onStartup bool) time.Time {
	if onStartup {
		return loop.reconcileUsersAndAccessPolicies()
	}

	for {
		// wait for next cycle of reconciliation if it's not the first "sync" on startup
		time.Sleep(time.Second * time.Duration(loop.reconciliationInterval))
		loop.reconcileUsersAndAccessPolicies()
	}
}

func (loop *ReconciliationLoop) reconcileUsersAndAccessPolicies() time.Time {
	logger := logs.GetDefaultLogger(loop.ctx)

	// 1. Get users from ACC with embedded access policies
	var accUsers map[string]*accclient.User
	var nextUpdateTimestamp time.Time
	err := retryHelper(loop.ctx,
		func() error {
			var getRequestError error
			accUsers, nextUpdateTimestamp, getRequestError = loop.getACCUsersAndAccessPoliciesForReconciliation(loop.ctx)
			return getRequestError
		})
	if err != nil {
		logger.Warnf("Failed to get ACC users: %v", err)
		return nextUpdateTimestamp // retry in next loop
	}

	// 2. Get users from External System
	var externalUserIDs map[string]struct{}
	err = retryHelper(loop.ctx,
		func() error {
			var getRequestError error
			externalUserIDs, getRequestError = loop.getExternalSystemUserIDs()
			return getRequestError
		})
	if err != nil {
		logger.Warnf("Failed to get external users: %v", err)
		return nextUpdateTimestamp // retry in next loop
	}

	// 3. remove non-existing users
	for externalUserID := range externalUserIDs {
		if _, ok := accUsers[externalUserID]; !ok {
			logger.Infof("Removing user %v", externalUserID)
			err = loop.extClient.DeleteUser(externalUserID)
			if err != nil {
				logger.Warnf("Failed to delete user %v: %v", externalUserID, err)
			}
		}
	}

	// 4. create or update users
	for accUserID, accUser := range accUsers {
		logger.Infof("Updating user %v", accUserID)
		err = createOrUpdateUser(loop.ctx, loop.extClient, loop.accClient, accUser)
		if err != nil {
			logger.Warnf("Failed to update user %v: %v", accUserID, err)
		}
	}

	// 5. get external access policies
	var externalAPs map[string]struct{}
	err = retryHelper(loop.ctx,
		func() error {
			var getRequestError error
			externalAPs, getRequestError = loop.getExternalSystemAccessPolicies()
			return getRequestError
		})
	if err != nil {
		logger.Warnf("Failed to get external access policies: %v", err)
		return nextUpdateTimestamp // retry in next loop
	}

	// 6. remove non existing access policies
	loop.deleteInactiveAccessPolicies(accUsers, externalAPs)

	// 7. create or update access policies
	loop.updateAccessPoliciesOnExternalSystem(accUsers)

	return nextUpdateTimestamp
}

// =====================
// helper functions
// =====================

// getACCTenantsAndOfferingItemsForReconciliation returns tenants information with embedded offering items that currently exist in ACC
// it doesn't use updated_since filter because we need current state of tenants and offering items for reconciliation purpose
func (loop *ReconciliationLoop) getACCTenantsAndOfferingItemsForReconciliation() (map[string]*accclient.Tenant, time.Time, error) {
	limit := uint(accPageSize)
	withContacts := true
	withOfferingItems := true
	tenantsRequest := &accclient.TenantGetRequest{
		SubTreeRootID:     loop.tenantID,
		Limit:             &limit,
		WithContacts:      &withContacts,
		WithOfferingItems: &withOfferingItems,
		// deleted tenants could be hard deleted by retention, still more accurate to pull tenants from ext-system
		AllowDeleted: false,
	}

	tenantsResp, err := loop.accClient.GetTenants(loop.ctx, tenantsRequest)
	if err != nil {
		return nil, time.Time{}, err
	}

	nextUpdateTimestamp := tenantsResp.Timestamp

	// create a map of tenantID to tenant object
	accTenants := make(map[string]*accclient.Tenant, len(tenantsResp.Items))
	for i := range tenantsResp.Items {
		if tenantsResp.Items[i].DeletedAt.IsZero() {
			accTenants[tenantsResp.Items[i].ID] = &tenantsResp.Items[i]
		}
	}

	if tenantsResp.After() == "" {
		// no second page
		return accTenants, nextUpdateTimestamp.Time, nil
	}

	// loop through the rest of the pages
	for tenantsResp != nil {
		tenantsResp, err = loop.accClient.GetTenantsNextPage(loop.ctx, tenantsResp)
		if err != nil {
			return nil, time.Time{}, err
		}
		if tenantsResp == nil {
			break
		}

		for i := range tenantsResp.Items {
			if tenantsResp.Items[i].DeletedAt.IsZero() {
				accTenants[tenantsResp.Items[i].ID] = &tenantsResp.Items[i]
			}
		}
	}

	return accTenants, nextUpdateTimestamp.Time, nil
}

// getExternalSystemTenantIDs returns a set of tenantIDs that currently exist in external system
func (loop *ReconciliationLoop) getExternalSystemTenantIDs() (map[string]struct{}, error) {
	tenantIDs := make(map[string]struct{})
	offset := 0
	for ; ; offset += externalSystemPageSize {
		tenants, err := loop.extClient.GetActiveTenantIDs(offset, externalSystemPageSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get tenants from external system: %v", err)
		}

		for i := range tenants {
			tenantIDs[tenants[i]] = struct{}{}
		}

		if len(tenants) < externalSystemPageSize {
			// last page
			break
		}
	}
	return tenantIDs, nil
}

// getExternalSystemOfferingItems returns a set of tenantIDs that currently exist in external system
func (loop *ReconciliationLoop) getExternalSystemOfferingItems() ([]core.OfferingItemID, error) {
	offeringItems := make([]core.OfferingItemID, 0)
	offset := 0
	for ; ; offset += externalSystemPageSize {
		externalOIs, err := loop.extClient.GetActiveOfferingItemIDs(offset, externalSystemPageSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get tenants from external system: %v", err)
		}

		offeringItems = append(offeringItems, externalOIs...)

		if len(externalOIs) < externalSystemPageSize {
			// last page
			break
		}
	}
	return offeringItems, nil
}

// deleteInactiveOfferingItems deletes offering item on external system if:
// 1. tenant already removed from ACC
// 2. tenant still exists in ACC but the particular offering item already disabled
// 3. tenant still exists in ACC but the particular offering item is no longer reported (already hard deleted by ACC)
func (loop *ReconciliationLoop) deleteInactiveOfferingItems(
	accTenants map[string]*accclient.Tenant,
	externalOIs []core.OfferingItemID) {
	logger := logs.GetDefaultLogger(loop.ctx)
	for i := range externalOIs {
		deleteOI := true
		if accOIs, ok := accTenants[externalOIs[i].TenantID]; ok {
			for j := range accOIs.OfferingItems {
				if accOIs.OfferingItems[j].Name == externalOIs[i].OfferingItemName {
					deleteOI = accOIs.OfferingItems[j].Status == 0
					break
				}
			}
		}

		if deleteOI {
			if err := loop.extClient.DeleteOfferingItem(externalOIs[i]); err != nil {
				logger.Warnf("Failed to delete offering item on external-system: %v", err)
			}
		}
	}
}

// updateOfferingItemsOnExternalSystem performs upsert for active offering items from ACC to external system
func (loop *ReconciliationLoop) updateOfferingItemsOnExternalSystem(accTenants map[string]*accclient.Tenant) {
	logger := logs.GetDefaultLogger(loop.ctx)
	for _, tenant := range accTenants {
		for i := range tenant.OfferingItems {
			if tenant.OfferingItems[i].Status > 0 {
				if oiCreated, err := loop.extClient.CreateOrUpdateOfferingItem(&tenant.OfferingItems[i]); err != nil {
					logger.Warnf("Failed to upsert offering item %v for tenant %v into external-system: %v",
						tenant.OfferingItems[i].Name, tenant.ID, err)
				} else {
					logger.Debugf("Offering item %v for tenant %v successfully updated (is new offering item: %v)",
						tenant.OfferingItems[i].Name, tenant.ID, oiCreated)
				}
			}
		}
	}
}

// getACCUsersAndAccessPoliciesForReconciliation returns user information with embedded access policies that currently exist in ACC
// it doesn't use updated_since filter because we need current state of users and access policies for reconciliation purpose
func (loop *ReconciliationLoop) getACCUsersAndAccessPoliciesForReconciliation(
	ctx context.Context) (map[string]*accclient.User, time.Time, error) {
	limit := uint(accPageSize)
	withAccessPolicies := true
	usersRequest := &accclient.UserGetRequest{
		SubTreeRootTenantID: loop.tenantID,
		WithAccessPolicies:  &withAccessPolicies,
		Limit:               &limit,
	}

	usersResp, err := loop.accClient.GetUsers(ctx, usersRequest)
	if err != nil {
		return nil, time.Time{}, err
	}

	nextUpdateTimestamp := usersResp.Timestamp

	// create a map of userID to user object
	accUsers := make(map[string]*accclient.User, len(usersResp.Items))
	for i := range usersResp.Items {
		if usersResp.Items[i].DeletedAt.IsZero() && usersResp.Items[i].ID != "" {
			accUsers[usersResp.Items[i].ID] = &usersResp.Items[i]
		}
	}

	if usersResp.After() == "" {
		// no second page
		return accUsers, nextUpdateTimestamp, nil
	}

	// loop through the rest of the pages
	for usersResp != nil {
		usersResp, err = loop.accClient.GetUsersNextPage(ctx, usersResp)
		if err != nil {
			return nil, time.Time{}, err
		}
		if usersResp == nil {
			break
		}

		for i := range usersResp.Items {
			if usersResp.Items[i].DeletedAt.IsZero() && usersResp.Items[i].ID != "" {
				accUsers[usersResp.Items[i].ID] = &usersResp.Items[i]
			}
		}
	}

	return accUsers, nextUpdateTimestamp, nil
}

// getExternalSystemUserIDs returns a set of userIDs that currently exist in external system
func (loop *ReconciliationLoop) getExternalSystemUserIDs() (map[string]struct{}, error) {
	userIDs := make(map[string]struct{})
	offset := 0
	for ; ; offset += externalSystemPageSize {
		userIDPage, err := loop.extClient.GetActiveUserIDs(offset, externalSystemPageSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get users from external system: %v", err)
		}

		for _, userID := range userIDPage {
			userIDs[userID] = struct{}{}
		}

		if len(userIDPage) < externalSystemPageSize {
			// last page
			break
		}
	}
	return userIDs, nil
}

// getExternalSystemAccessPolicies returns a set of policyIDs that currently exist in external system
func (loop *ReconciliationLoop) getExternalSystemAccessPolicies() (map[string]struct{}, error) {
	policyIDs := make(map[string]struct{})
	offset := 0
	for ; ; offset += externalSystemPageSize {
		policyIDPage, err := loop.extClient.GetActiveAccessPolicyIDs(offset, externalSystemPageSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get users from external system: %v", err)
		}

		for _, policyID := range policyIDPage {
			policyIDs[policyID] = struct{}{}
		}

		if len(policyIDPage) < externalSystemPageSize {
			// last page
			break
		}
	}
	return policyIDs, nil
}

// deleteInactiveAccessPolicies deletes access policy on external system if:
// 1. user already removed from ACC
// 2. user still exists in ACC but the particular access policy is soft deleted
// 3. user still exists in ACC but the particular access policy is no longer reported (already hard deleted by ACC)
func (loop *ReconciliationLoop) deleteInactiveAccessPolicies(
	accUsers map[string]*accclient.User,
	externalAPs map[string]struct{}) {
	logger := logs.GetDefaultLogger(loop.ctx)
	// externalAPs will be used as the set of policyIDs to delete
	// if externalAP ID matches accAP ID then it is removed from the map
	// and will not be deleted from externalSystem
	for i := range accUsers {
		for j := range accUsers[i].AccessPolicies {
			// skip deleted ACC policies
			if accUsers[i].AccessPolicies[j].DeletedAt != nil {
				continue
			}
			policyID := accUsers[i].AccessPolicies[j].ID
			// does nothing if key (policyID) doesnt exist in map (externalAPs)
			delete(externalAPs, policyID)
		}
	}

	for extPolicyID := range externalAPs {
		if err := loop.extClient.DeleteAccessPolicy(extPolicyID); err != nil {
			logger.Warnf("Failed to delete access policy on external-system: %v", err)
		}
	}
}

// updateAccessPoliciesOnExternalSystem performs upsert for active access policies from ACC to external system
func (loop *ReconciliationLoop) updateAccessPoliciesOnExternalSystem(accUsers map[string]*accclient.User) {
	logger := logs.GetDefaultLogger(loop.ctx)
	for _, user := range accUsers {
		for i := range user.AccessPolicies {
			if user.AccessPolicies[i].DeletedAt == nil {
				if apCreated, err := loop.extClient.CreateOrUpdateAccessPolicy(&user.AccessPolicies[i]); err != nil {
					logger.Warnf("Failed to upsert access policy %v with ID %v for user %v into external-system: %v",
						user.AccessPolicies[i].RoleID, user.AccessPolicies[i].ID, user.ID, err)
				} else {
					logger.Debugf("Access policy %v for user %v with ID %v successfully updated (is new access policy: %v)",
						user.AccessPolicies[i].RoleID, user.AccessPolicies[i].ID, user.ID, apCreated)
				}
			}
		}
	}
}
