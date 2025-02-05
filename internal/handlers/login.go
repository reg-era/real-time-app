package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"forum/internal/database"
	middleware "forum/internal/middleware"
	"forum/internal/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	type response struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var CredentialsUser response
	if err := json.NewDecoder(r.Body).Decode(&CredentialsUser); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}
	regexp := regexp.MustCompile(`^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$`)
	isemail := regexp.MatchString(CredentialsUser.Username)
	var userData utils.User

	if len(CredentialsUser.Username) < 5 || len(CredentialsUser.Password) < 8 || len(CredentialsUser.Username) > 30 || len(CredentialsUser.Password) > 64 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}
	if isemail {
		userData.Email = CredentialsUser.Username
		userData.Password = CredentialsUser.Password
	} else {
		userData.UserName = CredentialsUser.Username
		userData.Password = CredentialsUser.Password
	}

	password := userData.Password
	err := middleware.ValidCredential(db, &userData)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Incorect Username or password", http.StatusUnauthorized)
			return
		}
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	if !CheckPasswordHash(&password, &userData.Password) {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}
	ok, err := middleware.GetActiveSession(db, &userData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	if ok {
		err = middleware.DeleteSession(db, &userData)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			return
		}
	}

	userData.SessionId, err = GenerateSessionID()
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	userData.Expiration = time.Now().Add(1 * time.Hour)
	err = database.InsertSession(db, &userData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
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
