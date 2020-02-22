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

	"github.com/elastic/beats/libbeat/processors/dissect"
	"go.uber.org/zap"
)

// maxPostMemory memory limit for parsing the POST HTTP request
const maxPostMemory = 16 * 1024 * 1024

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
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
		tokenized := make([]map[string]string, len(samples))
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

	logger.Sugar().Infow("Server is running",
		"port", 8080,
	)

	http.ListenAndServe(":8080", nil)
}
