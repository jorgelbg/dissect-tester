// Copyright (c) 2020 Jorge Luis Betancourt. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAPIHandler(t *testing.T) {
	testCases := []struct {
		name           string
		path           string
		method         string
		payload        url.Values
		expectedBody   string
		expectedStatus int
	}{
		{
			name:   "Valid tokenizer and str",
			path:   "/api",
			method: "POST",
			payload: url.Values{
				"tokenizer": {"%{key1} %{key2}"},
				"str":       {"a b"},
			},
			expectedBody:   "[{\"key1\":\"a\",\"key2\":\"b\"}]",
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Missing str parameter",
			path:   "/api",
			method: "POST",
			payload: url.Values{
				"tokenizer": {"%{key1} %{key2}"},
			},
			expectedBody:   "str parameter not found\n",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, strings.NewReader(
				tc.payload.Encode()),
			)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(apiHandler)

			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			if diff := cmp.Diff(tc.expectedBody, rr.Body.String()); diff != "" {
				t.Errorf("handler returned wrong body (-want +got):\n%s", diff)
			}
		})
	}
}
