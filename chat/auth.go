package main

import (
	"log"
	"net/http"
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
