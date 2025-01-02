package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"forum/internal/utils"
)

func InsertOrUpdateReactionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userID int) {
	r.Header.Add("content-type", "application/json")

	reactionType := r.URL.Query().Get("reaction_type")
	targetType := r.URL.Query().Get("target_type")
	id := r.URL.Query().Get("target_id")

	if targetType != "" && reactionType != "" {
		var insertQuery string
		switch targetType {
		case "post":
			insertQuery = `INSERT INTO reactions (reaction_type, user_id, post_id, target_type) 
			VALUES (?, ?, ?, ?)
			ON CONFLICT (user_id, post_id, target_type) DO UPDATE SET reaction_type = EXCLUDED.reaction_type;
			`
		case "comment":
			insertQuery = `INSERT INTO reactions (reaction_type, user_id, comment_id, target_type) 
			VALUES (?, ?, ?, ?)
			ON CONFLICT (user_id, comment_id,target_type) DO UPDATE SET reaction_type = EXCLUDED.reaction_type ;
			`
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := db.Exec(insertQuery, reactionType, userID, id, targetType)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}
		w.WriteHeader(200)

	} else {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error": "Bad Request 1!"}`)
		return
	}
}

func DeleteReactionHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userID int) {
	r.Header.Add("content-type", "application/json")
	targetType := r.URL.Query().Get("target_type")
	id := r.URL.Query().Get("target_id")

	if targetType != "" && id != "" {
		var deleteQuery string
		switch targetType {
		case "post":
			deleteQuery = `DELETE FROM reactions WHERE user_id = ? AND post_id = ? AND target_type = ? `
		case "comment":
			deleteQuery = `DELETE FROM reactions WHERE user_id = ? AND comment_id = ? AND target_type = ? `
		}
		_, err := db.Exec(deleteQuery, userID, id, targetType)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
	} else {
		utils.RespondWithJSON(w, http.StatusBadRequest, `{"error": "Bad Request1"}`)
		return
	}
}

func GetReactionsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	targetID := r.URL.Query().Get("target_id")
	targetType := r.URL.Query().Get("target_type")

	var column string
	if targetType == "post" {
		column = "post_id"
	} else if targetType == "comment" {
		column = "comment_id"
	} else {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// Prepare query to get likes and dislikes for the target (post or comment)
	var likedBy, dislikedBy []int
	var userReaction string

	// Query for liked users
	likeQuery := fmt.Sprintf(`
		SELECT user_id
		FROM reactions
		WHERE %s = ? AND reaction_type = 'like' AND target_type = ? ;`, column)

	// Query for disliked users
	dislikeQuery := fmt.Sprintf(`
		SELECT user_id
		FROM reactions
		WHERE %s = ? AND reaction_type = 'dislike' AND target_type = ? ;`, column)

	// Query for user reaction to a post
	userReactionQuery := fmt.Sprintf(`
		SELECT reaction_type
		FROM reactions
		WHERE user_id = ? AND %s = ? AND target_type = ? ;`, column)

	// Execute like query
	rows, err := db.Query(likeQuery, targetID, targetType)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		likedBy = append(likedBy, userId)
	}

	// Execute dislike query
	rows, err = db.Query(dislikeQuery, targetID, targetType)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dislikedBy = append(dislikedBy, userId)
	}

	// Execute user reaction query
	err = db.QueryRow(userReactionQuery, userId, targetID, targetType).Scan(&userReaction)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := utils.Reaction{
		LikedBy:      likedBy,
		DislikedBy:   dislikedBy,
		UserReaction: userReaction,
	}

	// Send response
	utils.RespondWithJSON(w, http.StatusOK, response)
}
