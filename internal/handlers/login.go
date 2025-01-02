package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/internal/database"
	middleware "forum/internal/middleware"
	"forum/internal/utils"
	tmpl "forum/web"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, err := middleware.ValidUser(r, db)
	if err != nil {
		if err == http.ErrNoCookie {
			tmpl.ExecuteTemplate(w, []string{"login"}, http.StatusOK, nil)
			return
		}
		if err == sql.ErrNoRows {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Path:    "/",
				Value:   "",
				Expires: time.Unix(0, 0),
			})
			tmpl.ExecuteTemplate(w, []string{"login"}, http.StatusUnauthorized, nil)
			return
		}
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userData utils.User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	if len(userData.UserName) < 5 || len(userData.Password) < 8 || len(userData.UserName) > 30 || len(userData.Password) > 64 {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	password := userData.Password
	err := middleware.ValidCredential(db, &userData)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Incorect Username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !CheckPasswordHash(&password, &userData.Password) {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}
	ok, err := middleware.GetActiveSession(db, &userData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if ok {
		err = middleware.DeleteSession(db, &userData)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	userData.SessionId, err = GenerateSessionID()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userData.Expiration = time.Now().Add(1 * time.Hour)
	err = database.InsertSession(db, &userData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   userData.SessionId,
		Expires: userData.Expiration,
	})

	w.WriteHeader(http.StatusOK)
}
