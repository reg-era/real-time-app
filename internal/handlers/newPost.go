package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
)

func NewPostPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	categories, err := GetCategories(db)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		utils.RespondWithJSON(w, http.StatusOK, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, categories)
}

func NewPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fmt.Printf("Error parsing form: %v", err)
		utils.RespondWithJSON(w, http.StatusOK, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	post := &utils.Post{
		Title:     r.PostFormValue("title"),
		Content:   r.PostFormValue("content"),
		CreatedAt: time.Now(),
		UserId:    userId,
	}
	categories := r.Form["category"]
	if len(categories) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}

	if len(strings.TrimSpace(post.Title)) < 3 || len(strings.TrimSpace(post.Title)) > 60 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	} else if len(strings.TrimSpace(post.Content)) < 10 || len(strings.TrimSpace(post.Content)) > 10000 {
		utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
		return
	}

	_, err := database.InsertPost(post, db, categories)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
