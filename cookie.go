package main

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

func setSessionCookie(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("puki")
	if err != nil && err == http.ErrNoCookie {
		user_id_cookie := http.Cookie{
			Name:    "puki",
			Value:   uuid.New().String(),
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, &user_id_cookie)
	}
}
