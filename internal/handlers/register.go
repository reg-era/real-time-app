package handlers

import (
	"database/sql"
	"encoding/json"
	"html"
	"net/http"
	"time"

	"forum/internal/database"
	middleware "forum/internal/middleware"
	"forum/internal/utils"
	tmpl "forum/web"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	_, err := middleware.ValidUser(r, db)
	if err != nil {
		if err == http.ErrNoCookie {
			tmpl.ExecuteTemplate(w, []string{"register"}, http.StatusOK, nil)
			return
		}
		if err == sql.ErrNoRows {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Path:    "/",
				Value:   "",
				Expires: time.Unix(0, 0),
			})
			tmpl.ExecuteTemplate(w, []string{"register"}, http.StatusUnauthorized, nil)
			return
		}
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userData utils.User
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	if len(userData.UserName) < 5 || len(userData.Password) < 8 || len(userData.UserName) > 30 || len(userData.Password) > 64 || !isValidEmail(&userData.Email) {
		http.Error(w, "invalid username/password/email", http.StatusBadRequest)
		return
	}

	if userData.Password != userData.PasswordConfirmation {
		http.Error(w, "Passwords do not match.", http.StatusBadRequest)
		return
	}

	ok, err := middleware.IsUserRegistered(db, &userData)
	if err != nil {
		http.Error(w, "internaInternal Server Error", http.StatusInternalServerError)
		return
	}

	if ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	err = HashPassword(&userData.Password)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusNotAcceptable)
		return
	}

	userData.UserName = html.EscapeString(userData.UserName)
	err = middleware.RegisterUser(db, &userData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a session and set a cookie
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
