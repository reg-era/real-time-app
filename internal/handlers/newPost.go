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
	tmpl "forum/web"
)

func NewPostPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	categories, err := GetCategories(db)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, []string{"new_post", "sideBar"}, http.StatusOK, categories)
}

func NewPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if err := r.ParseForm(); err != nil {
		fmt.Printf("Error parsing form: %v", err)
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusBadRequest, http.StatusBadRequest)
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
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(post.Title)) < 3 || len(strings.TrimSpace(post.Title)) > 60 {
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusBadRequest, http.StatusBadRequest)
		return
	} else if len(strings.TrimSpace(post.Content)) < 10 || len(strings.TrimSpace(post.Content)) > 10000 {
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	_, err := database.InsertPost(post, db, categories)
	if err != nil {
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
