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

package accclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error contains the information about the error return from the server
type Error struct {
	Code    string                 `json:"code"`
	Context map[string]interface{} `json:"context,omitempty"`
	Domain  *string                `json:"domain"`
	Message string                 `json:"message"`
	Details Details                `json:"details"`
	Data    *[]string              `json:"data,omitempty"`
}

// Details holds additional information on the error.
type Details struct {
	Info *string `json:"info,omitempty"`
}

// Error implements the error interface
func (e Error) Error() string {
	return fmt.Sprintf("accclient error: domain: %v, reason: %v", getStrPtrValue(e.Domain), e.Code)
}

func getStrPtrValue(v *string) string {
	if v != nil {
		return *v
	}
	return nullString
}

func makeError(r *http.Response) error {
	if r.StatusCode >= 400 {
		return Error{
			Code: fmt.Sprintf("%d", r.StatusCode),
		}
	}

	// Other errors we decode the json error body
	var e Error
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return err
	}
	return &e
}
