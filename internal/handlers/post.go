package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	database "forum/internal/database"
	"forum/internal/utils"
	websocket "forum/internal/ws"
)

func PostsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int, ws *websocket.Hub) {
	id := r.URL.Query().Get("post_id")
	if id != "" {
		postId, err := strconv.Atoi(id)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
			return
		}

		post, err := database.ReadPost(db, userId, postId)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			return
		}

		post.Categories, err = database.GetPostCategories(db, post.PostId, userId)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, post)
		return
	}

	lastindex, err := database.GetLastPostId(db)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	json, err := json.Marshal(lastindex)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	_, err = w.Write(json)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
}
