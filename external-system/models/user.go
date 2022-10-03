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

package models

import (
	"errors"
	"time"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
)

const noItemError = "noItemError"

// User model for API response
type User struct {
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ID            string    `json:"id" gorm:"primaryKey"`
	TenantID      string    `json:"tenantId"`
	Login         string    `json:"login"`
	Contact       string    `json:"contact"`
	Activated     bool      `json:"isActivated"`
	Enabled       bool      `json:"isEnabled"`
	TermsAccepted bool      `json:"areTermsAccepted"`
}

// CreateOrUpdateUser creates user in db if it doesn't exist, updates if exist
func CreateOrUpdateUser(user *User) (bool, error) {
	var isCreated bool

	if _, err := GetUser(user.ID); err != nil {
		switch err.Error() {
		case noItemError:
			isCreated = true
		default:
			return false, err
		}
	}

	if err := config.DBConn.Save(&user).Error; err != nil {
		return false, err
	}

	return isCreated, nil
}

// GetUsers gets list of users from db
func GetUsers(offset, limit int) ([]User, error) {
	users := make([]User, 0, limit)
	if err := config.DBConn.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUser gets single user by id
func GetUser(id string) (User, error) {
	var user User
	if err := config.DBConn.Find(&user, "id = ?", id).Error; err != nil {
		return user, err
	}
	if user.ID != id {
		return user, errors.New(noItemError)
	}
	return user, nil
}

// DeleteUser soft deletes single user entry in db by id
func DeleteUser(id string) error {
	return config.DBConn.Where("id = ?", id).Delete(&User{}).Error
}
