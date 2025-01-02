package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	utils "forum/internal/utils"
	tmpl "forum/web"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func MeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	switch r.URL.Path {
	case "/me/liked_posts":
		query := `SELECT post_id FROM reactions WHERE user_id = ? AND reaction_type = 'like' AND post_id NOT NULL`
		rows, err := utils.QueryRows(db, query, userId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}

		postIds := []int{}
		for rows.Next() {
			var postId int
			err := rows.Scan(&postId)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
				return
			}
			postIds = append(postIds, postId)
		}

		jsonIds, err := json.Marshal(postIds)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, []string{"posts", "sideBar"}, http.StatusOK, string(jsonIds))
	case "/me/created_posts":
		query := `SELECT id FROM posts WHERE user_id = ?`
		rows, err := utils.QueryRows(db, query, userId)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}

		postIds := []int{}
		for rows.Next() {
			var postId int
			err := rows.Scan(&postId)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
			}
			postIds = append(postIds, postId)
		}

		jsonIds, err := json.Marshal(postIds)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusInternalServerError, http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, []string{"posts", "sideBar"}, http.StatusOK, string(jsonIds))
	default:
		tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusNotFound, http.StatusNotFound)
	}
}

func GenerateSessionID() (string, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return sessionID.String(), nil
}

func HashPassword(password *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	*password = string(bytes)
	return err
}

func CheckPasswordHash(password, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password))
	return err == nil
}

func isValidEmail(email *string) bool {
	*email = strings.ToLower(*email)
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(*email)
}
