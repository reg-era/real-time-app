package handlers

import (
	"database/sql"
	"net/http"

	"forum/internal/database"
	"forum/internal/utils"
	ws "forum/internal/ws"
)

func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int, hub *ws.Hub) {
	name := r.URL.Query().Get("name")
	data, err := database.GetConversations(db, userId, name)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, data)
}
