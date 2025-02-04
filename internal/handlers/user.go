package handlers

import (
	"database/sql"
	"net/http"
	"regexp"
	"strings"

	utils "forum/internal/utils"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func MeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	switch r.URL.Path {
	case "/api/me/liked-posts":
		query := `SELECT post_id FROM reactions WHERE user_id = ? AND reaction_type = 'like' AND post_id NOT NULL`
		rows, err := utils.QueryRows(db, query, userId)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			return
		}

		postIds := []int{}
		for rows.Next() {
			var postId int
			err := rows.Scan(&postId)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
				return
			}
			postIds = append(postIds, postId)
		}
		data := struct {
			PostIds []int `json:"post_ids"`
		}{PostIds: postIds}
		utils.RespondWithJSON(w, http.StatusOK, data)
	case "/api/me/created-posts":
		query := `SELECT id FROM posts WHERE user_id = ?`
		rows, err := utils.QueryRows(db, query, userId)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			return
		}

		postIds := []int{}
		for rows.Next() {
			var postId int
			err := rows.Scan(&postId)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			}
			postIds = append(postIds, postId)
		}

		data := struct {
			PostIds []int `json:"post_ids"`
		}{PostIds: postIds}
		utils.RespondWithJSON(w, http.StatusOK, data)

	case "/api/me/check-in":
		utils.RespondWithJSON(w, http.StatusAccepted, nil)

	default:
		utils.RespondWithJSON(w, http.StatusNotFound, utils.ErrorResponse{Error: "Status Not Found"})
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
