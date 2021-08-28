package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type templateHandler struct {
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := t.templ.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func newTemplateHandler(templatePath string) *templateHandler {
	t := &templateHandler{}
	t.templ = template.Must(template.ParseFiles(templatePath))
	return t
}

func main() {
	th := newTemplateHandler(filepath.FromSlash(`templates/chat.html`))
	http.Handle("/", th)

	r := newRoom()
	http.Handle("/room", r)

	go r.run()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
