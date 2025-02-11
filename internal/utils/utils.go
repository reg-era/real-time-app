package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

var Colors = map[string]string{"green": "\033[42m", "red": "\033[41m", "reset": "\033[0m"}

type User struct {
	UserId               int64     `json:"userId,omitempty"`
	UserName             string    `json:"username"`
	Email                string    `json:"email"`
	Password             string    `json:"password"`
	PasswordConfirmation string    `json:"confirmPassword"`
	Age                  string    `json:"age"`        // Changed to lowercase
	Gender               string    `json:"gender"`     // Changed to lowercase
	LastName             string    `json:"last_name"`  // Changed to match JSON
	FirstName            string    `json:"first_name"` // Changed to match JSON
	SessionId            string    `json:"sessionId,omitempty"`
	Expiration           time.Time `json:"expiration,omitempty"`
}

type Post struct {
	PostId     int
	UserId     int
	UserName   string
	Title      string
	Categories []string
	Content    string
	CreatedAt  time.Time
	Gender     string
}

type Reaction struct {
	LikedBy      []int  `json:"liked_by"`
	DislikedBy   []int  `json:"disliked_by"`
	UserReaction string `json:"user_reaction"`
}

type Message struct {
	Id         int       `json:"Id"`
	SenderID   int       `json:"SenderID"`
	ReceiverID int       `json:"ReceiverID"`
	Message    string    `json:"Message"`
	CreatedAt  time.Time `json:"CreatedAt"`
	IsSender   bool      `json:"IsSender"`
	SenderName string    `json:"sender_name"`
	Seen       int       `json:"Seen"`
}

type Comment struct {
	Comment_id int    `json:"comment_id"`
	Post_id    int    `json:"post_id"`
	User_id    int    `json:"user_id"`
	User_name  string `json:"user_name"`
	Content    string `json:"content"`
	Created_at string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func QueryRows(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func QueryRow(db *sql.DB, query string, args ...any) (*sql.Row, error) {
	stmt, err := db.Prepare(query)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer stmt.Close()
	return stmt.QueryRow(args...), nil
}
