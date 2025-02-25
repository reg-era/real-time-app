package handlers

import (
	"database/sql"
	"encoding/json"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"forum/internal/database"
	middleware "forum/internal/middleware"
	"forum/internal/utils"
	websocket "forum/internal/ws"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int, ws *websocket.Hub) {
	var userData utils.User

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Invalid input data"})
		return
	}

	regexp := regexp.MustCompile(`^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$`)
	isemail := regexp.MatchString(userData.Email)
	if !isemail {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "invalid email"})
		return
	}

	if age, err := strconv.Atoi(userData.Age); err != nil || age <= 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "invalid Age"})
		return
	}
	userData.UserName = strings.TrimSpace(userData.UserName)
	if len(userData.UserName) < 5 || len(userData.Password) < 8 || len(userData.UserName) > 30 || len(userData.Password) > 64 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "invalid username/password/email"})
		return
	}

	if userData.Password != userData.PasswordConfirmation {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Passwords do not match."})
		return
	}

	ok, err := middleware.IsUserRegistered(db, &userData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error "})
		return
	}

	if ok {
		utils.RespondWithJSON(w, http.StatusConflict, utils.ErrorResponse{Error: "User already exists"})
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
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	// Create a session and set a cookie
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
