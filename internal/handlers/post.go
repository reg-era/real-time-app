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
	tmpl "forum/web"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if r.URL.Path != "/" {
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusNotFound, http.StatusNotFound)
		return
	} else if r.Method != "GET" {
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusNotFound, http.StatusNotFound)
		return
	}

	tmpl.ExecuteTemplate(w, []string{"posts", "sideBar"}, http.StatusOK, nil)
}

func PostsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	id := r.URL.Query().Get("post_id")
	if id != "" {
		postId, err := strconv.Atoi(id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		post, err := database.ReadPost(db, userId, postId)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err)
			return
		}

		post.Categories, err = database.GetPostCategories(db, post.PostId, userId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, post)
		return
	}

	lastindex, err := database.GetLastPostId(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(lastindex)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(json)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
