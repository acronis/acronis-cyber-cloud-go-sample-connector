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
	"reflect"
	"strconv"
	"testing"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/server/handlers"
)

func generateSampleUser() models.User {
	return models.User{
		Activated:     true,
		Contact:       "contact",
		Enabled:       true,
		ID:            "id",
		Login:         "login",
		TenantID:      "tenantId",
		TermsAccepted: true,
	}
}

func TestClient_CreateOrUpdateUser(t *testing.T) {
	deleteAllUsers()
	client := NewClient(http.DefaultClient, testServerAddr)

	type args struct {
		user models.User
	}
	tests := []struct {
		name          string
		args          args
		wantIsCreated bool
		wantErr       bool
	}{
		{
			name: "it creates user successfully",
			args: args{
				user: generateSampleUser(),
			},
			wantIsCreated: true,
			wantErr:       false,
		},
		{
			name: "it updates user successfully",
			args: args{
				user: generateSampleUser(),
			},
			wantIsCreated: false,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateOrUpdateUser(&tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateOrUpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantIsCreated {
				t.Errorf("Client.CreateOrUpdateUser() = %v, wantIsCreated %v", got, tt.wantIsCreated)
			}
		})
	}
	deleteAllUsers()
}

func TestClient_DeleteUser(t *testing.T) {
	deleteAllUsers()
	client := NewClient(http.DefaultClient, testServerAddr)

	// Create user entry first before deletion
	user1 := generateSampleUser()
	if _, err := client.CreateOrUpdateUser(&user1); err != nil {
		t.Fatalf(err.Error())
	}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "it deletes user entry successfully",
			args: args{
				id: user1.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.DeleteUser(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	deleteAllUsers()
}

func TestClient_GetUsers(t *testing.T) {
	deleteAllUsers()
	client := NewClient(http.DefaultClient, testServerAddr)

	var userIds []string

	for i := 0; i < handlers.DefaultLimitValue+1; i++ {
		user := generateSampleUser()
		user.ID = strconv.Itoa(i)
		if isCreated, err := client.CreateOrUpdateUser(&user); err != nil || isCreated == false {
			t.Fatalf("test initialization fail for GetUsers")
		}
		userIds = append(userIds, user.ID)
	}

	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		args    args
		wantIds []string
		wantErr bool
	}{
		{
			name:    "it gets default number of users when limit is not provided with offset 0",
			args:    args{},
			wantIds: userIds[0:handlers.DefaultLimitValue],
			wantErr: false,
		},
		{
			name: "it gets correct number of users",
			args: args{
				offset: 0,
				limit:  2,
			},
			wantIds: userIds[0:2],
			wantErr: false,
		},
		{
			name: "it gets users at the correct offset",
			args: args{
				offset: 2,
				limit:  2,
			},
			wantIds: userIds[2:4],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetUsers(tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var gotIds []string
			for _, user := range got {
				gotIds = append(gotIds, user.ID)
			}
			if len(gotIds) != len(tt.wantIds) {
				t.Errorf("Client.GetUsers() = %v, want %v", len(gotIds), len(tt.wantIds))
			}
			if !reflect.DeepEqual(gotIds, tt.wantIds) {
				t.Errorf("Client.GetUsers() = %v, want %v", gotIds, tt.wantIds)
			}
		})
	}
	deleteAllUsers()
}

func TestClient_GetUser(t *testing.T) {
	deleteAllUsers()
	client := NewClient(http.DefaultClient, testServerAddr)

	// Create user entry first before deletion
	user1 := generateSampleUser()
	if _, err := client.CreateOrUpdateUser(&user1); err != nil {
		t.Fatalf(err.Error())
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantID  string
		wantErr bool
	}{
		{
			name: "it gets correct item",
			args: args{
				id: user1.ID,
			},
			wantID:  user1.ID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantID != got.ID {
				t.Errorf("Client.GetUser() = %v, want %v", got.ID, tt.wantID)
			}
		})
	}
	deleteAllUsers()
}

func deleteAllUsers() {
	config.DBConn.Unscoped().Exec("DELETE FROM users")
}
