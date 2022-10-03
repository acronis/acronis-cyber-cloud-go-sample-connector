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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is a external system client
type Client struct {
	APIURL     string
	HTTPClient *http.Client
}

// NewClient creates a new external system client
func NewClient(httpClient *http.Client, url string) *Client {
	return &Client{
		APIURL:     url,
		HTTPClient: httpClient,
	}
}

func (c *Client) do(method, url string, body interface{}) (*http.Response, error) {
	var reader io.Reader

	if body != nil {
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to create http request body. %w", err)
		}

		reader = bytes.NewBuffer(reqBody)
	} else {
		reader = nil
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request. %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server. %w", err)
	}

	return resp, nil
}

func (c *Client) doGet(endpoint string) (*http.Response, error) {
	url := c.APIURL + endpoint
	return c.do(http.MethodGet, url, nil)
}

func (c *Client) doPost(endpoint string, body interface{}) (*http.Response, error) {
	url := c.APIURL + endpoint
	return c.do(http.MethodPost, url, body)
}

func (c *Client) doDelete(endpoint string) (*http.Response, error) {
	url := c.APIURL + endpoint
	return c.do(http.MethodDelete, url, nil)
}
