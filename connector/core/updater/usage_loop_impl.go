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

// UsageLoop is a sample implementation that pulls usage information from external-system (ISV)
// and push these usage information into Acronis Cyber Cloud Platform
type UsageLoop struct {
	accClient *accclient.Client
	extClient core.ExternalSystemClient

	// optional to be set during initialization
	updateInterval uint // in seconds
}

// NewUsageLoop initializes UsageLoop as an implementation of core.UsageLoop
func NewUsageLoop(
	accClient *accclient.Client,
	extClient core.ExternalSystemClient,
	options ...func(*UsageLoop)) core.UsageLoop {
	loop := &UsageLoop{
		accClient:      accClient,
		extClient:      extClient,
		updateInterval: 21600, // default
	}

	for _, option := range options {
		option(loop)
	}

	return loop
}

// WithUsageUpdateInterval is an optional init function to set usage update interval
func WithUsageUpdateInterval(interval uint) func(*UsageLoop) {
	return func(loop *UsageLoop) {
		loop.updateInterval = interval
	}
}

// UpdateUsages will send usage report from external-system to ACC periodically
// 1. Get usages from external system
// 2. Push usage report to ACC
func (loop *UsageLoop) UpdateUsages() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.ContextID, "usage_loop")
	logger := logs.GetDefaultLogger(ctx)

	for ; ; time.Sleep(time.Second * time.Duration(loop.updateInterval)) {
		offset := 0
		for ; ; offset += externalSystemPageSize {
			// 1. Get usages from external-system
			pageUsages, err := loop.extClient.GetUsages(offset, externalSystemPageSize)
			if err != nil {
				// Retry whole loop if failed to get usage
				logger.Warnf("Failed to get external-system usages: %v", err)
				break
			}

			// No usages to send
			if len(pageUsages) == 0 {
				logger.Infof("No usages to push")
				break
			}

			logger.Infof("Pushing %v usages", len(pageUsages))
			// 2. Push usage report to ACC
			err = loop.sendACCUsageReport(ctx, pageUsages)
			if err != nil {
				// skip this batch if error
				logger.Warnf("Failed to push usages to ACC: %v", err)
			}

			if len(pageUsages) < externalSystemPageSize {
				// last page
				break
			}
		}
	}
}

// =====================
// helper functions
// =====================

func (loop *UsageLoop) sendACCUsageReport(ctx context.Context, extUsages []accclient.Usage) error {
	logger := logs.GetDefaultLogger(ctx)
	usageReq := &accclient.UsagesPutRequest{
		Items: extUsages,
	}

	usageResp, err := loop.accClient.UpdateUsages(ctx, usageReq)
	if err != nil {
		return err
	}

	logger.Infof("Successfully pushed %v usages", len(usageResp.Items))

	return nil
}
