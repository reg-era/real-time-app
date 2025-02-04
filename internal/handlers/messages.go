package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/internal/database"
	"forum/internal/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	name := r.URL.Query().Get("name")
	data, err := database.GetConversations(db, userId, name)
	fmt.Println(data)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, data)
}
