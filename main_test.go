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
		name     string
		method   string
		payload  url.Values
		response string
		status   int
	}{
		{
			name:   "Valid tokenizer and str",
			method: "POST",
			payload: url.Values{
				"tokenizer": {"%{key1} %{key2}"},
				"str":       {"a b"},
			},
			response: "[{\"key1\":\"a\",\"key2\":\"b\"}]",
			status:   http.StatusOK,
		},
		{
			name:   "Missing str parameter",
			method: "POST",
			payload: url.Values{
				"tokenizer": {"%{key1} %{key2}"},
			},
			response: "str parameter not found\n",
			status:   http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "/api", strings.NewReader(
				tc.payload.Encode()),
			)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type",
				"application/x-www-form-urlencoded;charset=UTF-8")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(apiHandler)

			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.status)
			}

			if diff := cmp.Diff(tc.response, rr.Body.String()); diff != "" {
				t.Errorf("handler returned wrong body (-want +got):\n%s", diff)
			}
		})
	}
}
