package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/elastic/beats/libbeat/processors/dissect"
)

func main() {
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

		fmt.Printf("str=%+v, tokenizer=%+v\n", str, tokenizer)
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

	println("server is running on : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
