package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/internal/database"
	"forum/internal/utils"
	ws "forum/internal/ws"
)

func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int, hub *ws.Hub) {
	name := r.URL.Query().Get("name")
	_, err := database.GetUserIdByName(name, db)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	err = database.Updatesenn(userId, db)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	data, err := database.GetConversations(db, userId, name)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, data)
}
