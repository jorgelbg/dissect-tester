// Copyright (c) 2020 Jorge Luis Betancourt. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/beats/v7/libbeat/processors/dissect"
	"github.com/elastic/beats/v7/libbeat/version"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// several timeout options for the HTTP server
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
	procTimeout  = 3 * time.Second
)

// A list of HTTP endpoints to register
const (
	staticPath = "/static/"
	apiPath    = "/api/"
)

var versionInfo = struct {
	Version string
}{
	Version: version.GetDefaultVersion(),
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	if err := tmpl.Execute(w, versionInfo); err != nil {
		zap.L().Error("Could not parse template.",
			zap.String("template", "templates/index.html"),
			zap.Error(err),
		)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't parse POST request: %s", err.Error()),
			http.StatusBadRequest)
		return
	}

	str := r.Form.Get("str")
	if len(str) == 0 {
		http.Error(w, "str parameter not found", http.StatusBadRequest)
		return
	}

	tokenizer := r.Form.Get("tokenizer")
	if len(tokenizer) == 0 {
		http.Error(w, "tokenizer parameter not found", http.StatusBadRequest)
		return
	}

	zap.L().Sugar().Infow("Received request",
		"str", str,
		"tokenizer", tokenizer,
	)

	samples := strings.Split(str, "\n")
	tokenized := make([]map[string]interface{}, 0)
	for i, s := range samples {
		processor, err := dissect.New(tokenizer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		m, err := processor.DissectConvert(s)
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
	w.Write(payload) // nolint: errcheck
}

// RegisterAppHandlers registers the app handlers with the given mux
func RegisterAppHandlers(mux *http.ServeMux) {
	mux.Handle(staticPath,
		http.StripPrefix(staticPath, http.FileServer(http.Dir("static"))),
	)

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc(apiPath, apiHandler)
}

func main() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		panic("Couldn't configure the logger. Aborting!")
	}

	zap.ReplaceGlobals(logger)
	defer logger.Sync() // nolint: errcheck

	logger.Info("elastic/beats engine", zap.String("version", versionInfo.Version))

	listenAddr := flag.String("listen", ":8080", "address to listen on")
	flag.Parse()

	_, err = maxprocs.Set(maxprocs.Logger(
		func(logMessage string, args ...interface{}) {
			logger.Sugar().Info(fmt.Sprintf(logMessage, args...))
		},
	))

	if err != nil {
		logger.Error("Failed to set maxprocs: %v", zap.Error(err))
	}

	mux := http.NewServeMux()
	pprofMux := http.NewServeMux()

	RegisterAppHandlers(mux)
	RegisterDebugHandler(pprofMux)

	go func() {
		pprofServer := http.Server{
			Addr:    "localhost:6060",
			Handler: pprofMux,
		}
		logger.Info("Starting debug server")
		logger.Error("Error starting pprof server", zap.Error(pprofServer.ListenAndServe()))
	}()

	server := http.Server{
		Handler:      http.TimeoutHandler(mux, procTimeout, "Processing your request took too long."),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	defer server.Close()
	server.Addr = *listenAddr

	logger.Sugar().Infow("Server is running",
		"address", *listenAddr,
	)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("Could not start HTTP server.",
			zap.Error(err),
		)
	}
}
