package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"json_parser/lexer"
	"json_parser/parser"
)

type pageData struct {
	Input  string
	Output string
	Errors []string
}

var pageTmpl = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>JSON Parser Demo</title>
  <style>
    body { font-family: sans-serif; margin: 2rem auto; max-width: 900px; padding: 0 1rem; }
    textarea { width: 100%; min-height: 220px; font-family: monospace; font-size: 14px; }
    button { margin-top: 0.8rem; padding: 0.5rem 1rem; }
    pre { background: #f6f8fa; padding: 1rem; overflow: auto; }
    .errors { color: #b00020; }
  </style>
</head>
<body>
  <h1>JSON Parser</h1>
  <form method="post" action="/parse">
    <textarea name="json">{{.Input}}</textarea>
    <br>
    <button type="submit">Parse</button>
  </form>

  {{if .Errors}}
    <h2>Errors</h2>
    <ul class="errors">
      {{range .Errors}}<li>{{.}}</li>{{end}}
    </ul>
  {{end}}

  {{if .Output}}
    <h2>Parsed Output</h2>
    <pre>{{.Output}}</pre>
  {{end}}
</body>
</html>`))

func renderPage(w http.ResponseWriter, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTmpl.Execute(w, data); err != nil {
		http.Error(w, "template render error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, pageData{Input: `{"name":"Shreyash","age":25,"active":true,"tags":["go","json"]}`})
	})

	http.HandleFunc("/parse", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		input := r.FormValue("json")
		l := lexer.New(input)
		p := parser.New(l)
		result := p.Parse()

		data := pageData{Input: input}
		if errs := p.Errors(); len(errs) > 0 {
			data.Errors = errs
		} else {
			data.Output = fmt.Sprintf("%#v", result)
		}

		renderPage(w, data)
	})

	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
