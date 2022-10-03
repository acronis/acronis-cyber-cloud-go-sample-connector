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

// Reconciliation is an interface to perform reconciliation logic
// between Acronis Cyber Cloud Platform and external system database.
type Reconciliation interface {
	// ReconcileTenantsAndOfferingItems will reconcile tenants and offering items objects.
	// If onStartup is set to true, it will only run once and return the timestamp that can be used
	// for the next "update loop"
	// If onStartup is set to false, it will run periodically every ReconciliationInterval set in config file
	ReconcileTenantsAndOfferingItems(onStartup bool) time.Time

	// ReconcileUsersAndAccessPolicies will reconcile users and access policies objects.
	// If onStartup is set to true, it will only run once and return the timestamp that can be used
	// for the next "update loop"
	// If onStartup is set to false, it will run periodically every ReconciliationInterval set in config file
	ReconcileUsersAndAccessPolicies(onStartup bool) time.Time
}
