package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/hiroygo/go-programming-blueprints/trace"
)

type templateHandler struct {
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := t.templ.Execute(w, r); err != nil {
		log.Println(err)
	}
}

func newTemplateHandler(templatePath string) *templateHandler {
	t := &templateHandler{}
	t.templ = template.Must(template.ParseFiles(templatePath))
	return t
}

func parseArgs() string {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	return *addr
}

func main() {
	addr := parseArgs()

	th := newTemplateHandler(filepath.FromSlash(`templates/chat.html`))
	http.Handle("/", th)

	// r := newRoom(trace.New(os.Stderr))
	r := newRoom(trace.New(nil))
	http.Handle("/room", r)

	go r.run()
	log.Printf("サーバ開始: %q\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
