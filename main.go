// Copyright (c) 2020 Jorge Luis Betancourt. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/processors/dissect"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// maxPostMemory memory limit for parsing the POST HTTP request
const maxPostMemory = 16 * 1024 * 1024

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
)

// A list of HTTP endpoints to register
const (
	StaticPath = "/static/"
	APIPath    = "/api/"
)

func main() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	defer logger.Sync()

	_, err := maxprocs.Set(maxprocs.Logger(
		func(logMessage string, args ...interface{}) {
			logger.Sugar().Info(fmt.Sprintf(logMessage, args...))
		},
	))

	if err != nil {
		logger.Error("Failed to set maxprocs: %v", zap.Error(err))
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})

	mux.Handle(StaticPath, http.StripPrefix(StaticPath, http.FileServer(http.Dir("static"))))

	mux.HandleFunc(APIPath, func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(maxPostMemory)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't parse POST request: %s", err.Error()),
				http.StatusBadRequest)
			return
		}

		str, _ := url.QueryUnescape(r.Form.Get("str"))
		if len(str) == 0 {
			http.Error(w, "samples parameter not found", http.StatusBadRequest)
			return
		}

		tokenizer, _ := url.QueryUnescape(r.Form.Get("tokenizer"))
		if len(tokenizer) == 0 {
			http.Error(w, "pattern parameter not found", http.StatusBadRequest)
			return
		}

		logger.Sugar().Infow("Received request",
			"str", str,
			"tokenizer", tokenizer,
		)

		samples := strings.Split(str, "\n")
		tokenized := make([]map[string]string, 0)
		for i, s := range samples {
			processor, err := dissect.New(tokenizer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			m, err := processor.Dissect(s)
			if err != nil {
				http.Error(w, fmt.Sprintf("sample: %d, error: %s", i, err), http.StatusBadRequest)
				return
			}

			tokenized = append(tokenized, m)
		}

		payload, err := json.Marshal(tokenized)
		if err != nil {
			http.Error(w, "couldn't encode response", http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	})

	RegisterDebugHandler(mux)

	server := http.Server{
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	server.Addr = ":8080"

	logger.Sugar().Infow("Server is running",
		"port", 8080,
	)

	server.ListenAndServe()
}
