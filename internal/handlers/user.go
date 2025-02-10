package handlers

import (
	"database/sql"
	"net/http"

	"forum/internal/database"
	utils "forum/internal/utils"

	websocket "forum/internal/ws"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func MeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int, websocket *websocket.Hub) {
	if r.URL.Path == "/api/me/check-in" {
		name, err := database.GetUserName(userId, db)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, nil)
			return
		}

		utils.RespondWithJSON(w, http.StatusAccepted, name)

	} else {
		utils.RespondWithJSON(w, http.StatusNotFound, utils.ErrorResponse{Error: "Status Not Found"})
	}
}

func GenerateSessionID() (string, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return sessionID.String(), nil
}

func HashPassword(password *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	*password = string(bytes)
	return err
}

func CheckPasswordHash(password, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password))
	return err == nil
}
