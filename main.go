package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/elastic/beats/libbeat/processors/dissect"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		str, ok := r.URL.Query()["str"]
		if !ok || len(str[0]) == 0 {
			http.Error(w, "str query parameter not found", http.StatusBadRequest)
		}

		tokenizer, ok := r.URL.Query()["tokenizer"]
		if !ok || len(tokenizer[0]) < 1 {
			http.Error(w, "tokenizer query parameter not found", http.StatusBadRequest)
			return
		}

		logger.Sugar().Infow("Received request",
			"str", str[0],
			"tokenizer", tokenizer[0],
		)

		processor, err := dissect.New(tokenizer[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		m, err := processor.Dissect(str[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payload, err := json.Marshal(m)
		if err != nil {
			http.Error(w, "couldn't encode response", http.StatusNoContent)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	})

	logger.Sugar().Infow("Server is running",
		"port", 8080,
	)
	http.ListenAndServe(":8080", nil)
}
