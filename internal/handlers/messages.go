package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"forum/internal/database"
	"forum/internal/utils"
)

func GetConversations(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	conversations, err := database.GetConversations(db, userId)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, conversations)
}

func PostMessage(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	var message utils.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}

	message.SenderID = userId
	err := database.CreateMessage(&message, db)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, message)
}

// db *sql.DB, userId, limit, from int
func GetMessages(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	query, err := getQuerys(r)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}

	messages, err := database.GetMessages(db, userId, query[0], query[1])
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, messages)
}

func getQuerys(r *http.Request) ([2]int, error) {
	res := [2]int{}
	var err error
	for i, key := range []string{"limit", "from"} {
		value := r.URL.Query().Get(key)
		if value == "" {
			return [2]int{}, errors.New("missing " + key)
		}
		res[i], err = strconv.Atoi(value)
		if err != nil {
			return [2]int{}, errors.New("failed to get " + key + " value")
		}
	}
	return res, nil
}
