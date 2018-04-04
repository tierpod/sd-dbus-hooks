package main

import (
	"log"
	"net/http"
)

type tokenStore struct {
	token string
}

func (t *tokenStore) middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Token")
		if t.token == "" || t.token == token {
			log.Printf("[INFO] xtoken: accept request from %v to %v: X-Token \"%v\"", r.RemoteAddr, r.URL.Path, token)
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("[ERROR] xtoken: forbidden request from %v to %v: wrong X-Token header: \"%v\"", r.RemoteAddr, r.URL.Path, token)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	return http.HandlerFunc(fn)
}

func newTokenStore(token string) *tokenStore {
	return &tokenStore{
		token: token,
	}
}
