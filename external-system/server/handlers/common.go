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

// The `handlers` package provides handler functions for routes
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var (
	// DefaultLimitValue is default number of items returned is limit is not specified in requests
	DefaultLimitValue = 10
)

// API response for error scenarios
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Write API Response in JSON format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		response, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, "JSON encoding failed:"+err.Error(), http.StatusInternalServerError)
			return
		}

		_, writeErr := w.Write(response)
		if writeErr != nil {
			log.Printf("Write failed: %v", writeErr)
		}
	}
}

func getOptionalIntQueryParam(r *http.Request, varKey string) (int, error) {
	valueStr, ok := r.URL.Query()[varKey]
	if !ok {
		return 0, nil
	}

	value, err := strconv.Atoi(valueStr[0])
	if err != nil {
		return 0, err
	}

	return value, nil
}
