package auth

import (
	"database/sql"
	"net/http"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
)

type customHandler func(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int)

func AuthMiddleware(db *sql.DB, next customHandler, login bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isConstentJson := r.Header.Get("Content-Type") == "application/json"
		userId, err := ValidUser(r, db)
		if err != nil {
			if err == http.ErrNoCookie {
				if isConstentJson {
					utils.RespondWithJSON(w, http.StatusUnauthorized, `{"error":"Unauthorized"}`)
					return
				}
				if login {
					next(w, r, db, userId)
					return
				}
				utils.RespondWithJSON(w, http.StatusUnauthorized, utils.ErrorResponse{Error: "Unauthorized"})
				return
			} else if err == sql.ErrNoRows {
				http.SetCookie(w, &http.Cookie{
					Name:    "session_token",
					Path:    "/",
					Value:   "",
					Expires: time.Unix(0, 0),
				})
				if isConstentJson {
					utils.RespondWithJSON(w, http.StatusUnauthorized, `{"error":"Unauthorized"}`)
					return
				}
				if login {
					next(w, r, db, userId)
					return
				}
				utils.RespondWithJSON(w, http.StatusUnauthorized, utils.ErrorResponse{Error: "Unauthorized"})
				return
			} else {
				utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
				return
			}
		}
		next(w, r, db, userId)
	})
}

func IsUserRegistered(db *sql.DB, userData *utils.User) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?);`
	err := db.QueryRow(query, userData.Email, userData.UserName).Scan(&exists)
	return exists, err
}

func RegisterUser(db *sql.DB, userData *utils.User) error {
	insertQuery := `INSERT INTO users (username, Age, Gender, First_Name, Last_Name, email, password) VALUES (?, ?, ?, ?, ?, ?, ?);`
	result, err := db.Exec(insertQuery, userData.UserName, userData.Age, userData.Gender, userData.FirstName, userData.LastName, userData.Email, userData.Password)
	if err != nil {
		return err
	}
	userData.UserId, err = result.LastInsertId()
	return err
}

func GetActiveSession(db *sql.DB, userData *utils.User) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM sessions WHERE user_id = ?  AND expires_at > ?);`
	err := db.QueryRow(query, userData.UserId, userData.Expiration).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func DeleteSession(db *sql.DB, userData *utils.User) error {
	query := `DELETE FROM sessions WHERE user_id =  ?;`
	_, err := db.Exec(query, userData.UserId)
	return err
}

func ValidCredential(db *sql.DB, userData *utils.User) error {
	query := `SELECT id, password FROM users WHERE username = ?;`
	err := db.QueryRow(query, userData.UserName).Scan(&userData.UserId, &userData.Password)
	if err != nil {
		return err
	}
	return err
}

func ValidUser(r *http.Request, db *sql.DB) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}
	userid, err := database.Get_session(cookie.Value, db)
	if err != nil {
		return 0, err
	}
	return userid, nil
}

func RemoveUser(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cookie.Value)
	if err != nil {
		return err
	}
	return nil
}
