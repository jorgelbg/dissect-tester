// Copyright (c) 2020 Jorge Luis Betancourt. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"net/http"
	"net/http/pprof"
	"strings"
)

const (
	pprofPath = "/debug/pprof/"
)

// RegisterDebugHandler registers the pprof handler with the given http mux.
func RegisterDebugHandler(mux *http.ServeMux) {
	mux.Handle(pprofPath, handler())
}

func handler() http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, pprofPath)
		switch name {
		case "cmdline":
			pprof.Cmdline(w, r)
		case "profile":
			pprof.Profile(w, r)
		case "symbol":
			pprof.Symbol(w, r)
		case "trace":
			pprof.Trace(w, r)
		default:
			pprof.Index(w, r)
		}
	}

	return http.HandlerFunc(h)
}
