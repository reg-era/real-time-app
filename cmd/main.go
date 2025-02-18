package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"forum/internal/database"
	"forum/internal/handlers"
	auth "forum/internal/middleware"
	"forum/internal/utils"
	websocket "forum/internal/ws"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	port := os.Getenv("PORT")

	db := database.CreateDatabase(dbPath)
	forumHub := &websocket.Hub{
		Clients:    make(map[int][]*websocket.Client),
		Broadcast:  make(chan *sql.DB),
		Register:   make(chan *websocket.Client),
		Message:    make(chan utils.Message),
		Unregister: make(chan *websocket.Client),
		Progress:   make(chan websocket.Progresser),
	}
	go forumHub.Run()
	defer close(forumHub.Broadcast)
	defer close(forumHub.Register)
	defer close(forumHub.Message)
	defer close(forumHub.Unregister)
	defer close(forumHub.Progress)
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
		auth.AuthMiddleware(db, func(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int, hub *websocket.Hub) {
			http.ServeFile(w, r, "web/index.html")
		}, 100, time.Minute, nil, false).ServeHTTP(w, r)
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
			auth.AuthMiddleware(db, handlers.RegisterHandler, 1, 30*time.Second, nil, false).ServeHTTP(w, r)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			auth.AuthMiddleware(db, handlers.LoginHandler, 10, 30*time.Second, nil, false).ServeHTTP(w, r)
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
		if r.Method == "GET" {
			auth.AuthMiddleware(db, handlers.PostsHandler, 300, time.Minute, nil, true).ServeHTTP(w, r)
		} else {
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
			return
		}
	})

	router.HandleFunc("/api/new_post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.NewPostPageHandler, 200, time.Minute, nil, true).ServeHTTP(w, r)
		case "POST":
			auth.AuthMiddleware(db, handlers.NewPostHandler, 3, time.Minute, nil, true).ServeHTTP(w, r)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.GetCommentsHandler, 200, time.Minute, nil, true).ServeHTTP(w, r)
		case "POST":
			auth.AuthMiddleware(db, handlers.AddCommentHandler, 10, time.Minute, nil, true).ServeHTTP(w, r)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/api/react", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "GET" {
			auth.AuthMiddleware(db, handlers.GetReactionsHandler, 100, time.Second, nil, true).ServeHTTP(w, r)
		} else if method == "PUT" {
			auth.AuthMiddleware(db, handlers.InsertOrUpdateReactionHandler, 10, time.Second, nil, true).ServeHTTP(w, r)
		} else if method == "DELETE" {
			auth.AuthMiddleware(db, handlers.DeleteReactionHandler, 10, time.Second, nil, true).ServeHTTP(w, r)
		} else {
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
			return
		}
	})

	router.HandleFunc("/api/me/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.MeHandler, 10, time.Second, nil, true).ServeHTTP(w, r)
		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware(db, handlers.HandleWs, 1000, time.Minute, forumHub, true).ServeHTTP(w, r)
	})

	router.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			auth.AuthMiddleware(db, handlers.GetUser, 100, time.Second, nil, true).ServeHTTP(w, r)

		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	log.Printf("Route server running on http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}
