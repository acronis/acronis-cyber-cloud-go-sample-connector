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

import "time"

// SyncLoop is an interface to perform sync operations with pull-based mechanism.
// It's intended to capture incremental updates from Acronis Cyber Cloud Platform and
// push these information into external-system
type SyncLoop interface {
	// UpdateTenantsAndOfferingItems captures updates of tenants and offering items
	// firstUpdatedSince can be supplied with timestamp returned by the first reconciliation upon startup
	UpdateTenantsAndOfferingItems(firstUpdatedSince time.Time)

	// UpdateUsersAndAccessPolicies captures updates of useres and access policies
	// firstUpdatedSince can be supplied with timestamp returned by the first reconciliation upon startup
	UpdateUsersAndAccessPolicies(firstUpdatedSince time.Time)
}
