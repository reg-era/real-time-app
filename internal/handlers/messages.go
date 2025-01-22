package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	name := r.URL.Query().Get("name")
	user, err := database.GetUserIdByName(name, db)
	if err != nil {
		if user == -69 {
			utils.RespondWithJSON(w, http.StatusNotFound, utils.ErrorResponse{Error: "User not found"})
			return
		}
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil)
}

func GetAllFriends(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	data, err := database.GetAllFriends(db, userId)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	// check if users online or not befor sending data
	utils.RespondWithJSON(w, http.StatusOK, data)
}

func GetConversations(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	name := r.URL.Query().Get("name")

	data, err := database.GetConversations(db, userId, name)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, data)
}

func PostMessage(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	respons := struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&respons); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}

	var message utils.Message
	var err error
	message.SenderID = userId
	message.Message = respons.Message
	message.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	if message.ReceiverID, err = database.GetUserIdByName(respons.To, db); err != nil {
		if message.ReceiverID == -69 {
			utils.RespondWithJSON(w, http.StatusNotFound, utils.ErrorResponse{Error: "User not found"})
			return
		}
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
	}

	err = database.CreateMessage(&message, db)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, message)
}
