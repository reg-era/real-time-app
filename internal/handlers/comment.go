package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
)

type querys struct {
	postid int
	limit  int
	from   int
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request, file *sql.DB, userId int) {
	var data querys
	err := getDataQuery(&data, r)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: fmt.Sprintf(`{"error":"status bad request %v"}`, err)})
		return
	}

	comments, err := database.GetComments(data.postid, file, userId, data.limit, data.from)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, comments)
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int, pool *utils.Pool) {
	userName, err := database.GetUserName(userId, db)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	comment := utils.Comment{}
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	comment.Content = strings.TrimSpace(comment.Content)
	if len(strings.TrimSpace(comment.Content)) < 1 || len(comment.Content) > 2000 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Comment must be between 3 and 2000 characters"})
		return
	}

	comment.User_name = userName
	comment.User_id = userId
	comment.Created_at = time.Now().Format(time.RFC3339)

	if err := database.CreateComment(&comment, db); err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, comment)
}

func getDataQuery(field *querys, r *http.Request) error {
	allKeys := []string{"post", "from", "limit"}
	for _, key := range allKeys {
		data, err := strconv.Atoi(r.URL.Query().Get(key))
		if err != nil {
			return errors.New("faild to get " + key + " value")
		}
		switch key {
		case "post":
			field.postid = data
		case "from":
			field.from = data
		case "limit":
			field.limit = data
		}
	}
	return nil
}
