package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	database "forum/internal/database/models"
	utils "forum/internal/utils"
)

func CreateDatabase(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("%sError opening database:%s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	}
	// Verify the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("%sError accessing database: %s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	} else {
		fmt.Printf("%sDatabase created/opened successfully!%s\n", utils.Colors["green"], utils.Colors["reset"])
	}

	_, err = db.Exec(`PRAGMA foreign_keys=ON;`)
	if err != nil {
		log.Fatalf("%sError enabling foreign keys: %s%s\n", utils.Colors["red"], err.Error(), utils.Colors["reset"])
	}
	return db
}

func CreateTables(db *sql.DB) {
	_, err := db.Exec(database.UsersTable + database.SessionsTable + database.ReactionTable + database.MessageTable +
		database.CommentsTable + database.PostsTable + database.CategoriesTable + database.PostCategoriesTable)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Created all tables succesfully")
}

func CleanupExpiredSessions(db *sql.DB) {
	_, err := db.Exec("DELETE FROM sessions WHERE  expires_at < ?", time.Now())
	if err != nil {
		fmt.Printf("Error cleaning up expired sessions: %v", err)
	}
}

func InsertPost(p *utils.Post, db *sql.DB, categories []string) (int64, error) {
	transaction, err := db.Begin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting transaction:", err)
		return 0, err
	}
	stmt, err := transaction.Prepare(`INSERT INTO posts(user_id ,title,content) Values (?,?,?);`)
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "Error Adding post:", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.UserId, p.Title, p.Content)
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "Error Adding post:", err)
		return 0, err
	}
	lastPostID, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "error in assigning category to post", err)
		return 0, err
	}

	err = LinkPostWithCategory(transaction, categories, lastPostID, p.UserId)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		fmt.Fprintln(os.Stderr, "transaction aborted")
		return 0, err
	}
	return lastPostID, nil
}

func ReadPost(db *sql.DB, userId int, postId int) (*utils.Post, error) {
	query := `SELECT * FROM posts WHERE id = ?`
	row, err := utils.QueryRow(db, query, postId)
	if err != nil {
		return nil, err
	}

	Post := &utils.Post{}
	err = row.Scan(&Post.PostId, &Post.UserId, &Post.Title, &Post.Content, &Post.CreatedAt)
	if err != nil {
		return nil, err
	}

	Post.UserName, err = GetUserName(int(Post.UserId), db)
	if err != nil {
		return nil, err
	}
	return Post, nil
}

func GetLastPostId(db *sql.DB) (int, error) {
	query := `SELECT COALESCE(MAX(id), 0) FROM posts`
	row, err := utils.QueryRow(db, query)
	if err != nil {
		return 0, err
	}

	result := 0
	err = row.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func Get_session(ses string, db *sql.DB) (int, error) {
	var sessionid int
	query := `SELECT user_id FROM sessions WHERE session_id = ?`
	err := db.QueryRow(query, ses).Scan(&sessionid)
	if err != nil {
		return 0, err
	}
	return sessionid, nil
}

func InsertSession(db *sql.DB, userData *utils.User) error {
	_, err := db.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", userData.SessionId, userData.UserId, userData.Expiration)
	return err
}

func CreateComment(c *utils.Comment, db *sql.DB) error {
	query := `
	INSERT INTO comments (user_id, post_id, content, created_at)
	VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(query, c.User_id, c.Post_id, c.Content, c.Created_at)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.Comment_id = int(id)
	return nil
}

func GetComments(postID int, db *sql.DB, userId, limit, from int) ([]utils.Comment, error) {
	query := `
	SELECT comments.id, comments.content, comments.created_at, users.username  FROM comments
	INNER JOIN users ON comments.user_id = users.id
	WHERE comments.post_id = ?
	ORDER BY comments.created_at ASC
	LIMIT ? OFFSET ?;
	`
	rows, err := utils.QueryRows(db, query, postID, limit, from)
	if err != nil {
		return nil, errors.New(err.Error() + "here 1")
	}
	defer rows.Close()

	var comments []utils.Comment
	for rows.Next() {
		var comment utils.Comment
		err := rows.Scan(&comment.Comment_id, &comment.Content, &comment.Created_at, &comment.User_name)
		if err != nil {
			return nil, errors.New(err.Error() + "here 2")
		}
		comment.Post_id = postID
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New(err.Error() + "here 3")
	}

	return comments, nil
}

func GetCategoryContentIds(db *sql.DB, categoryId string) ([]int, error) {
	rows, err := utils.QueryRows(db, "SELECT post_id FROM post_categories WHERE category_id=?", categoryId)
	if err != nil {
		return nil, err
	}

	var ids []int
	for rows.Next() {
		tmp := 0
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		ids = append(ids, tmp)
	}

	return ids, nil
}

func GetUserName(id int, db *sql.DB) (string, error) {
	var name string
	row, err := utils.QueryRow(db, "SELECT username FROM users WHERE id = ?", id)
	if err != nil {
		return "", err
	}

	err = row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetPostCategories(db *sql.DB, PostId int, userId int) ([]string, error) {
	query := `
	SELECT categories.name 
	FROM post_categories
	JOIN categories ON categories.id = post_categories.category_id
	AND post_categories.post_id = ?;
	`

	rows, err := utils.QueryRows(db, query, PostId)
	if err != nil {
		return nil, err
	}

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func LinkPostWithCategory(transaction *sql.Tx, categories []string, postId int64, userId int) error {
	for _, category := range categories {
		stmt, err := transaction.Prepare(`INSERT INTO post_categories(user_id, post_id, category_id) VALUES(?, ?, ?);`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		tmp, err := strconv.Atoi(category)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(userId, postId, tmp)
		if err != nil {
			return err
		}
	}
	return nil
}

///////// only for messages ////////////

func GetConversations(db *sql.DB, userId int) ([]string, error) {
	query := `
	SELECT DISTINCT users.username
	FROM messages
	INNER JOIN users ON messages.sender_id = users.id
	WHERE messages.receiver_id = ?
	`
	rows, err := utils.QueryRows(db, query, userId)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	var conversations []string
	for rows.Next() {
		var conversation string
		err := rows.Scan(&conversation)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		conversations = append(conversations, conversation)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New(err.Error())
	}

	return conversations, nil
}

func CreateMessage(m *utils.Message, db *sql.DB) error {
	query := `
	INSERT INTO messages (sender_id, receiver_id, content, created_at)
	VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(query, m.SenderID, m.ReceiverID, m.Message, m.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.SenderID = int(id)
	return nil
}

func GetMessages(db *sql.DB, userId, limit, from int) ([]utils.Message, error) {
	query := `
	SELECT messages.id, messages.content, messages.created_at, users.username  FROM messages
	INNER JOIN users ON messages.sender_id = users.id
	WHERE messages.receiver_id = ?
	ORDER BY messages.created_at ASC
	LIMIT ? OFFSET ?;
	`
	rows, err := utils.QueryRows(db, query, userId, limit, from)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	var messages []utils.Message
	for rows.Next() {
		var message utils.Message
		err := rows.Scan(&message.SenderID, &message.Message, &message.CreatedAt, &message.ReceiverID)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		message.ReceiverID = userId
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New(err.Error())
	}

	return messages, nil
}
