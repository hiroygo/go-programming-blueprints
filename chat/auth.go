package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == nil {
		a.next.ServeHTTP(w, r)
		return
	}
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	log.Fatalf("ServeHTTP error: %v", err)
}

func MustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ss := strings.Split(r.URL.Path, "/")
	if len(ss) != 4 {
		log.Fatalf("len(Split(%v)) returns %v\n", r.URL.Path, len(ss))
	}

	action := ss[2]
	provider := ss[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalf("Provider error: %v\n", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalf("GetBeginAuthURL error: %v\n", err)
		}

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalf("Provider error: %v\n", err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalf("CompleteAuth error: %v\n", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalf("GetUser error: %v\n", err)
		}

		// JSON 形式の文字列を base64 変換する
		// e.g. `{"name":"yamada taro"}`
		authValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authValue,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %q not supported", action)
	}
}
