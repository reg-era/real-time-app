package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"forum/internal/database"
	"forum/internal/handlers"
	auth "forum/internal/middleware"
	"forum/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	port := os.Getenv("PORT")

	db := database.CreateDatabase(dbPath)
	forumHub := new(utils.Hub)
	defer db.Close()

	database.CreateTables(db)

	go func() {
		for {
			database.CleanupExpiredSessions(db)
			time.Sleep(2 * time.Hour)
		}
	}()

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	router.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, "/error?status=404", http.StatusMovedPermanently)
			return
		}
		fs := http.FileServer(http.Dir("web/assets/"))
		http.StripPrefix("/api/", fs).ServeHTTP(w, r)
	})

	router.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handlers.RegisterHandler(w, r, db)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handlers.LoginHandler(w, r, db)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
			return
		}
		auth.RemoveUser(w, r, db)
	})

	router.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
			return
		}
		userId, _ := auth.ValidUser(r, db)
		handlers.PostsHandler(w, r, db, userId)
	})

	router.HandleFunc("/api/new_post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.NewPostPageHandler, false).ServeHTTP(w, r)
		case "POST":
			auth.AuthMiddleware(db, handlers.NewPostHandler, false).ServeHTTP(w, r)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			userId, _ := auth.ValidUser(r, db)
			handlers.GetCommentsHandler(w, r, db, userId)
		case "POST":
			auth.AuthMiddleware(db, handlers.AddCommentHandler, false).ServeHTTP(w, r)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/api/react", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "GET" {
			userId, _ := auth.ValidUser(r, db)
			handlers.GetReactionsHandler(w, r, db, userId)
		} else if method == "PUT" {
			auth.AuthMiddleware(db, handlers.InsertOrUpdateReactionHandler, false).ServeHTTP(w, r)
		} else if method == "DELETE" {
			auth.AuthMiddleware(db, handlers.DeleteReactionHandler, false).ServeHTTP(w, r)
		} else {
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
			return
		}
	})

	router.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.CategoriesHandler(w, r, db, 0)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/api/me/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.MeHandler, false).ServeHTTP(w, r)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// userId, err := auth.ValidUser(r, db)
		// if err != nil {
		// 	utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		// }
		handlers.HandleWs(w, r, 1, db, forumHub)
	})

	// router.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case "GET":
	// 		section := r.URL.Query().Get("section")
	// 		switch section {
	// 		case "user":
	// 			name := r.URL.Query().Get("name")
	// 			if name == "" {
	// 				auth.AuthMiddleware(db, handlers.GetAllFriends, false).ServeHTTP(w, r)
	// 				return
	// 			}
	// 			auth.AuthMiddleware(db, handlers.GetUser, false).ServeHTTP(w, r)
	// 		case "message":
	// 			auth.AuthMiddleware(db, handlers.GetConversations, false).ServeHTTP(w, r)
	// 			return
	// 		default:
	// 			utils.RespondWithJSON(w, http.StatusBadRequest, utils.ErrorResponse{Error: "Bad Request"})
	// 		}
	// 	case "POST":
	// 		auth.AuthMiddleware(db, handlers.PostMessage, false).ServeHTTP(w, r)
	// 	default:
	// 		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
	// 	}
	// })

	log.Printf("Route server running on http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}
