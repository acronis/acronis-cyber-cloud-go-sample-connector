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
	"net/http"
	"strconv"
	"testing"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

type testReportingType int

const (
	perResource testReportingType = iota
	perTenant
	perTenantInfra
)

func testGenerateUsage(id uint, reportingType testReportingType) models.Usage {
	usage := models.Usage{
		ID:         id,
		UsageValue: int64(id),
	}
	switch reportingType {
	case perResource:
		rid := "resourceID" + strconv.Itoa(int(id))
		usageType := "usageType"
		usage.ResourceID = &rid
		usage.UsageType = &usageType

	case perTenant:
		offeringItem := "OfferingItem" + strconv.Itoa(int(id))
		usage.OfferingItem = &offeringItem

	case perTenantInfra:
		var infraID = "infra id"
		offeringItem := "OfferingItem" + strconv.Itoa(int(id))
		usage.OfferingItem = &offeringItem
		usage.InfraID = &infraID
	}
	return usage
}

func testPointerStrCompare(want, got *string) bool {
	return (want == nil && got == nil) || (*want == *got)
}

func testCompareUsageDto(want, got models.Usage) bool {
	rid := testPointerStrCompare(want.ResourceID, got.ResourceID)
	ut := testPointerStrCompare(want.UsageType, got.UsageType)
	tid := testPointerStrCompare(want.TenantID, got.TenantID)
	of := testPointerStrCompare(want.OfferingItem, got.OfferingItem)
	iid := testPointerStrCompare(want.InfraID, got.InfraID)
	uv := want.UsageValue == got.UsageValue

	return rid || ut || tid || of || iid || uv
}

func TestClient_GetUsages(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	config.DBConn.Unscoped().Exec("DELETE FROM usages")

	// init items
	var usages = []models.Usage{
		testGenerateUsage(1, perTenant),
		testGenerateUsage(2, perTenantInfra),
		testGenerateUsage(3, perResource),
		testGenerateUsage(4, perResource),
		testGenerateUsage(5, perResource),
		testGenerateUsage(6, perResource),
		testGenerateUsage(7, perResource),
		testGenerateUsage(8, perResource),
		testGenerateUsage(9, perResource),
		testGenerateUsage(10, perResource),
		testGenerateUsage(11, perResource),
	}
	config.DBConn.Create(&usages)

	// want values list
	var usageList []models.Usage
	for i, usage := range usages {
		switch i {
		case 0:
			usageList = append(usageList, usage)
		case 1:
			usageList = append(usageList, usage)
		default:
			usageList = append(usageList, usage)
		}
	}

	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Usage
		wantErr bool
	}{
		{
			name: "it gets default number of usages when limit is not provided, at offset 0",
			args: args{},
			want: usageList[0:10],
		},
		{
			name: "it gets correct number of usages",
			args: args{
				offset: 0,
				limit:  3,
			},
			want: usageList[0:3],
		},
		{
			name: "it gets usages at correct offset",
			args: args{
				offset: 2,
				limit:  2,
			},
			want: usageList[2:4],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetUsages(tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetUsages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Wrong number of items, got = %v, want %v", len(got), len(tt.want))
			}
			// compare items
			for i, gotUsageDto := range got {
				if wantUsageDto := tt.want[i]; !testCompareUsageDto(wantUsageDto, gotUsageDto) {
					t.Errorf("Client.GetUsages() items[%v] = %v, want %v", i, gotUsageDto, wantUsageDto)
				}
			}
		})
	}
}

func TestClient_CreateOrUpdateUsage(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	config.DBConn.Unscoped().Exec("DELETE FROM usages")

	usage1 := testGenerateUsage(1, perTenant)

	type args struct {
		usage *models.Usage
	}
	tests := []struct {
		name          string
		args          args
		wantIsCreated bool
		wantErr       bool
	}{
		{
			name: "it creates usage successfully",
			args: args{
				usage: &usage1,
			},
			wantIsCreated: true,
			wantErr:       false,
		},
		{
			name: "it updates usage successfully",
			args: args{
				usage: &usage1,
			},
			wantIsCreated: false,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateOrUpdateUsage(tt.args.usage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateOrUpdateUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantIsCreated {
				t.Errorf("Client.CreateOrUpdateUsage() = %v, want %v", got, tt.wantIsCreated)
			}
		})
	}
}
