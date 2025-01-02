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
	tmpl "forum/web"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	port := os.Getenv("PORT")

	// Create Database file
	db := database.CreateDatabase(dbPath)
	defer db.Close()

	// Create tables if not exist
	database.CreateTables(db)

	go func() {
		for {
			database.CleanupExpiredSessions(db)
			time.Sleep(2 * time.Hour)
		}
	}()

	// Create a multipluxer
	router := http.NewServeMux()

	router.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusNotFound, http.StatusNotFound)
			return
		}
		fs := http.FileServer(http.Dir("web/assets"))
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	// HomePage handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		auth.AuthMiddleware(db, handlers.HomePageHandler, true).ServeHTTP(w, r)
	})

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.RegisterPageHandler, true).ServeHTTP(w, r)
		case "POST":
			handlers.RegisterHandler(w, r, db)
		default:
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.LoginPageHandler(w, r, db)
		case "POST":
			handlers.LoginHandler(w, r, db)
		default:
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
		auth.RemoveUser(w, r, db)
	})

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
		userId, _ := auth.ValidUser(r, db)
		handlers.PostsHandler(w, r, db, userId)
	})

	router.HandleFunc("/new_post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.NewPostPageHandler, false).ServeHTTP(w, r)
		case "POST":
			auth.AuthMiddleware(db, handlers.NewPostHandler, false).ServeHTTP(w, r)
		default:
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			userId, _ := auth.ValidUser(r, db)
			handlers.GetCommentsHandler(w, r, db, userId)
		case "POST":
			auth.AuthMiddleware(db, handlers.AddCommentHandler, false).ServeHTTP(w, r)
		default:
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/react", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "GET" {
			userId, _ := auth.ValidUser(r, db)
			handlers.GetReactionsHandler(w, r, db, userId)
		} else if method == "PUT" {
			auth.AuthMiddleware(db, handlers.InsertOrUpdateReactionHandler, false).ServeHTTP(w, r)
		} else if method == "DELETE" {
			auth.AuthMiddleware(db, handlers.DeleteReactionHandler, false).ServeHTTP(w, r)
		} else {
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})

	router.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.CategoriesHandler(w, r, db, 0)
		default:
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/me/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			auth.AuthMiddleware(db, handlers.MeHandler, false).ServeHTTP(w, r)
		default:
			tmpl.ExecuteTemplate(w, []string{"error"}, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Route server running on http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, router))
}
