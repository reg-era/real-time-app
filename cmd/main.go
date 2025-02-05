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
		Clients:    make(map[*websocket.Client]int),
		Broadcast:  make(chan *sql.DB),
		Register:   make(chan *websocket.Client),
		Message:    make(chan utils.Message),
		Logout:     make(chan *websocket.Client),
		Unregister: make(chan *websocket.Client),
	}
	go forumHub.Run()
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
			auth.AuthMiddleware(db, handlers.GetReactionsHandler, false).ServeHTTP(w, r)
		} else if method == "PUT" {
			auth.AuthMiddleware(db, handlers.InsertOrUpdateReactionHandler, false).ServeHTTP(w, r)
		} else if method == "DELETE" {
			auth.AuthMiddleware(db, handlers.DeleteReactionHandler, false).ServeHTTP(w, r)
		} else {
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
			return
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
		userId, err := auth.ValidUser(r, db)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
			return
		}
		handlers.HandleWs(w, r, userId, db, forumHub)
	})

	router.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			auth.AuthMiddleware(db, handlers.GetUser, false).ServeHTTP(w, r)

		default:
			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, utils.ErrorResponse{Error: "Status Method Not Allowed"})
		}
	})

	log.Printf("Route server running on http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}
