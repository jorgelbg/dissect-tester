// Copyright (c) 2020 Jorge Luis Betancourt. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/elastic/beats/v7/libbeat/processors/dissect"
	"github.com/google/go-cmp/cmp"
	"github.com/magiconair/properties/assert"
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
		{
			name:   "Missing tokenizer parameter",
			method: "POST",
			payload: url.Values{
				"str": {"a b"},
			},
			response: "tokenizer parameter not found\n",
			status:   http.StatusBadRequest,
		},
		{
			name:     "Empty payload",
			method:   "GET",
			payload:  nil,
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

func TestIndexRoute(t *testing.T) {
	mux := http.NewServeMux()
	RegisterAppHandlers(mux)

	// The response recorder used to record HTTP responses
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Creating 'GET /' request failed,")
	}

	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("'GET /' request should've a handler registered, got: %d, want: %d",
			rr.Code, http.StatusOK,
		)
	}
}

func TestAPIRoute(t *testing.T) {
	mux := http.NewServeMux()
	RegisterAppHandlers(mux)

	// The response recorder used to record HTTP responses
	rr := httptest.NewRecorder()
	payload := url.Values{
		"tokenizer": {"%{key1} %{key2}"},
		"str":       {"a b"},
	}

	req, err := http.NewRequest("POST", APIPath, strings.NewReader(
		payload.Encode()))
	if err != nil {
		t.Fatalf("Creating 'POST /api' request failed.")
	}
	req.Header.Add("Content-Type",
		"application/x-www-form-urlencoded;charset=UTF-8")

	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("'POST /api' request should've a handler registered, got: %d, want: %d",
			rr.Code, http.StatusOK,
		)
	}
}

func TestBadUserCase(t *testing.T) {
	const (
		pattern = "%{id} %{function-\u003e}%{server}"
		message = `00000043 ViewReceive machine-321    `
	)

	_, err := dissect.New(pattern)
	myError := errors.New("invalid dissect tokenizer")
	assert.Equal(t, myError, err)
}
