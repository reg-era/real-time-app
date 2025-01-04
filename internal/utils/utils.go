package utils

import (
	"time"
)

var Colors = map[string]string{"green": "\033[42m", "red": "\033[41m", "reset": "\033[0m"}

type User struct {
	UserId               int64
	Nickname             string `json:"nickname"`
	Age                  int    `json:"age"`
	Gender               string `json:"gender"`
	Firstname            string `json:"firstname"`
	Lastname             string `json:"lastname"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"confirmPassword"`
	SessionId            string
	Expiration           time.Time
}

type Post struct {
	PostId     int
	UserId     int
	UserName   string
	Title      string
	Categories []string
	Content    string
	CreatedAt  time.Time
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
