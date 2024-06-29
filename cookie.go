package main

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

func setSessionCookie(w http.ResponseWriter, r *http.Request) http.Cookie {
	session_id, session_err := r.Cookie("chikiyt_session_id")

	if session_err != nil {
		if session_err == http.ErrNoCookie {
			session_id = &http.Cookie{
				Name:     "chikiyt_session_id",
				Value:    uuid.New().String(),
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, session_id)
		}
	}

	return *session_id

	// cookie, err := r.Cookie("chikiyt_session_id")
	//
	// if err != nil && err == http.ErrNoCookie {
	// 	cookie := http.Cookie{
	// 		Name:    "chikiyt_session_id",
	// 		Value:   uuid.New().String(),
	// 		Expires: time.Now().Add(24 * time.Hour),
	// 	}
	// 	http.SetCookie(w, &cookie)
	// 	return cookie
	// } else {
	// 	return *cookie
	// }
}

//
// func readSessionCookieValue(r *http.Request) string {
// 	cookie, err := r.Cookie("chikiyt_session_id")
// 	if err != nil && err == http.ErrNoCookie {
// 		return ""
//
// 	}
// 	return cookie.Value
// }
