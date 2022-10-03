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

const defaultMaxRetries = 5 // number of times to retry failed requests

// getTenantsByUUIDs is a helper function to GetTenants from Acronis Cyber Cloud Platform with UUIDs filter
// It requires accClient to interact with Acronis Cloud and list of tenantUUIDs to be retrieved.
func getTenantsByUUIDs(ctx context.Context, accClient *accclient.Client, tenantUUIDs []string) (*accclient.TenantGetResponse, error) {
	getTenantReq := &accclient.TenantGetRequest{
		UUIDs: tenantUUIDs,
	}
	return accClient.GetTenants(ctx, getTenantReq)
}

// createOrUpdateTenant is a helper function to enable recursively creating/updating tenants on external-system
// It requires ExternalSystem client implementation to interact with external system,
// accClient to interact with Acronis Cloud, and tenant object to be created/updated.
func createOrUpdateTenant(ctx context.Context,
	extClient core.ExternalSystemClient,
	accClient *accclient.Client,
	tenant *accclient.Tenant) error {
	logger := logs.GetDefaultLogger(ctx)

	// check if parent tenant exists only if it's non-root
	// root tenant has tenant.ID == tenant.ParentID, simply create this tenant in this case
	if tenant.ParentID != tenant.ID {
		if parentExists, err := extClient.CheckTenantExist(tenant.ParentID); err != nil {
			return fmt.Errorf("failed to check tenant existence for %v from external-system: %w", tenant.ParentID, err)
		} else if !parentExists {
			// if parent tenant not found, try to create parent first
			parentTenantResp, err := getTenantsByUUIDs(ctx, accClient, []string{tenant.ParentID})
			if err != nil {
				return fmt.Errorf("failed to get tenant %v from ACC: %w", tenant.ParentID, err)
			}
			if len(parentTenantResp.Items) == 0 {
				return fmt.Errorf("empty tenant response for %v", tenant.ParentID)
			}

			// recursively try to create parent tenant
			if err := createOrUpdateTenant(ctx, extClient, accClient, &parentTenantResp.Items[0]); err != nil {
				return fmt.Errorf("failed to create parent tenant with ID %v: %w", parentTenantResp.Items[0], err)
			}
		}
	}

	tenantCreated, err := extClient.CreateOrUpdateTenant(tenant)
	if err != nil {
		logger.Warnf("Failed to upsert tenant %v: %v", tenant.ID, err)
		return fmt.Errorf("failed to upsert tenant with ID %v: %w", tenant.ID, err)
	}
	logger.Debugf("Tenant %v successfully updated (is new tenant: %v)", tenant.ID, tenantCreated)

	return nil
}

// createOrUpdateUser is a helper function which will call the recursive function
// to create tenants if it does not exist for the given user and then creates/updates user
func createOrUpdateUser(ctx context.Context,
	extClient core.ExternalSystemClient,
	accClient *accclient.Client,
	user *accclient.User) error {
	logger := logs.GetDefaultLogger(ctx)

	// check if tenant exists for user
	if tenantExists, err := extClient.CheckTenantExist(user.TenantID); err != nil {
		return fmt.Errorf("failed to check tenant existence for %v from external-system: %w", user.TenantID, err)
	} else if !tenantExists {
		// if tenant not found, get the tenant item and try to create it
		tenantResp, err := getTenantsByUUIDs(ctx, accClient, []string{user.TenantID})
		if err != nil {
			return fmt.Errorf("failed to get tenant %v from ACC: %w", user.TenantID, err)
		}
		if len(tenantResp.Items) == 0 {
			return fmt.Errorf("empty tenant response for %v", user.TenantID)
		}

		// use recursive function to create tenants
		if err := createOrUpdateTenant(ctx, extClient, accClient, &tenantResp.Items[0]); err != nil {
			return fmt.Errorf("failed to create tenant with ID %v: %w", tenantResp.Items[0], err)
		}
	}

	userCreated, err := extClient.CreateOrUpdateUser(user)
	if err != nil {
		logger.Warnf("Failed to upsert user %v: %v", user.ID, err)
		return fmt.Errorf("failed to update user %v: %v", user.ID, err)
	}
	logger.Debugf("User %v successfully updated (is new user: %v)", user.ID, userCreated)

	return nil
}

// retryHelper is a helper function that retries the passed in function up to max retries on error,
// backing off exponential amount of time after each try
func retryHelper(ctx context.Context, userFunction func() error) error {
	var err error
	logger := logs.GetDefaultLogger(ctx)
	initialBackOff := 1 * time.Second
	for i := 0; i < defaultMaxRetries; i++ {
		err = userFunction()
		if err != nil {
			if i+1 >= defaultMaxRetries {
				break
			}
			logger.Warnf("failed retry %v: %v", i+1, err)
			time.Sleep(initialBackOff)
			initialBackOff *= 2
			continue
		}
		return nil
	}
	return fmt.Errorf("max %v retries reached: %v", defaultMaxRetries, err)
}
