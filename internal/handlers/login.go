package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/internal/database"
	middleware "forum/internal/middleware"
	"forum/internal/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var userData utils.User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}
	if len(userData.UserName) < 5 || len(userData.Password) < 8 || len(userData.UserName) > 30 || len(userData.Password) > 64 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
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
