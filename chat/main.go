package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hiroygo/go-programming-blueprints/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
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
	clientId := os.Getenv("AUTH_GOOGLE_ID")
	if clientId == "" {
		log.Fatal("AUTH_GOOGLE_ID is empty")
	}
	clientSecret := os.Getenv("AUTH_GOOGLE_SECRET")
	if clientSecret == "" {
		log.Fatal("AUTH_GOOGLE_SECRET is empty")
	}

	// p.45
	// クライアントとサーバ間で処理の進行状況をやり取りする際にデジタル署名を行う
	// デジタル署名により、データ改ざんを防げる
	gomniauth.SetSecurityKey("mysecretkey")
	gomniauth.WithProviders(
		google.New(clientId, clientSecret, "http://localhost:8080/auth/callback/google"),
	)

	chat := newTemplateHandler(filepath.FromSlash(`templates/chat.html`))
	http.Handle("/chat", MustAuth(chat))

	login := newTemplateHandler(filepath.FromSlash(`templates/login.html`))
	http.Handle("/login", login)

	http.HandleFunc("/auth/", loginHandler)

	// r := newRoom(trace.New(os.Stderr))
	r := newRoom(trace.New(nil))
	http.Handle("/room", r)

	go r.run()
	log.Printf("サーバ開始: %q\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
